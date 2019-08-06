package shared

import (
	"archive/zip"
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lestrrat/go-jwx/jwa"
	"github.com/lestrrat/go-jwx/jwt"
	"github.com/spf13/viper"
	types "github.com/srinandan/apigeecli/cmd/types"
)

const BaseURL = "https://apigee.googleapis.com/v1/organizations/"

var RootArgs = types.Arguments{}

//log levels, default is error
var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

//LogInfo controls the log level
var LogInfo = false

//skip checking access token expiry
var skipCheck = true

//skip writing access token to file
var skipCache = false

// Structure to hold OAuth response
type oAuthAccessToken struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
}

//EntityPayloadList stores list of entities
var EntityPayloadList [][]byte //types.EntityPayloadList

const accessTokenFile = ".access_token"

//Init function initializes the logger objects
func Init() {

	var infoHandle = ioutil.Discard

	if LogInfo {
		infoHandle = os.Stdout
	}

	warningHandle := os.Stdout
	errorHandle := os.Stdout

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

//This method is used to send resources, proxy bundles, shared flows etc.
func PostHttpOctet(print bool, url string, proxyName string) (respBody []byte, err error) {

	file, _ := os.Open(proxyName)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("proxy", proxyName)
	if err != nil {
		Error.Fatalln("Error writing multi-part: ", err)
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		Error.Fatalln("error copying multi-part: ", err)
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		Error.Fatalln("error closing multi-part: ", err)
		return nil, err
	}
	client := &http.Client{}

	Info.Println("Connecting to : ", url)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		Error.Fatalln("error in client: ", err)
		return nil, err
	}

	Info.Println("Setting token : ", RootArgs.Token)
	req.Header.Add("Authorization", "Bearer "+RootArgs.Token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)

	if err != nil {
		Error.Fatalln("error connecting: ", err)
		return nil, err
	} else {
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Error.Fatalln("error in response: ", err)
			return nil, err
		} else if resp.StatusCode != 200 {
			Error.Fatalln("error in response: ", string(respBody))
			return nil, errors.New("Error in response")
		}
		if print {
			return respBody, PrettyPrint(respBody)
		}
		return respBody, nil
	}
}

//This method is used to download resources, proxy bundles, sharedflows
func DownloadResource(url string, name string) error {

	out, err := os.Create(name + ".zip")
	if err != nil {
		Error.Fatalln("error creating file: ", err)
		return err
	}
	defer out.Close()

	client := &http.Client{}

	Info.Println("Connecting to : ", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Error.Fatalln("error in client: ", err)
		return err
	}

	Info.Println("Setting token : ", RootArgs.Token)
	req.Header.Add("Authorization", "Bearer "+RootArgs.Token)

	resp, err := client.Do(req)

	if err != nil {
		Error.Fatalln("error connecting: ", err)
		return err
	} else if resp.StatusCode > 299 {
		Error.Fatalln("response Code: ", resp.StatusCode)
		Error.Fatalln("error in response: ", resp.Body)
		return errors.New("error in response")
	} else {
		defer resp.Body.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			Error.Fatalln("error writing response to file: ", err)
			return err
		}

		fmt.Println("Proxy bundle " + name + ".zip completed")
		return nil
	}
}

// The first parameter instructs whether the output should be printed
// The second parameter is url. If only one parameter is sent, assume GET
// The third parameter is the payload. The two parameters are sent, assume POST
// THe fourth parammeter is the method. If three parameters are sent, assume method in param
func HttpClient(print bool, params ...string) (respBody []byte, err error) {

	var req *http.Request

	client := &http.Client{}
	Info.Println("Connecting to: ", params[0])

	if len(params) > 3 {
		return nil, errors.New("Incorrect parameters to invoke the method")
	}

	if len(params) == 1 {
		req, err = http.NewRequest("GET", params[0], nil)
	} else if len(params) == 2 {
		Info.Println("Payload: ", params[1])
		req, err = http.NewRequest("POST", params[0], bytes.NewBuffer([]byte(params[1])))
	} else if len(params) == 3 {
		if params[2] == "DELETE" {
			req, err = http.NewRequest("DELETE", params[0], nil)
		} else if params[2] == "PUT" {
			req, err = http.NewRequest("PUT", params[0], bytes.NewBuffer([]byte(params[1])))
		} else {
			return nil, errors.New("Unsupported method")
		}
	}

	if err != nil {
		Error.Fatalln("error in client: ", err)
		return nil, err
	}

	Info.Println("Setting token : ", RootArgs.Token)
	req.Header.Add("Authorization", "Bearer "+RootArgs.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		Error.Fatalln("error connecting: ", err)
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		Error.Fatalln("error in response: ", err)
		return nil, err
	} else if resp.StatusCode > 299 {
		Error.Fatalln("response Code: ", resp.StatusCode)
		Error.Fatalln("error in response: ", string(respBody))
		return nil, errors.New("Error in response")
	}
	if print {
		return respBody, PrettyPrint(respBody)
	}
	return respBody, nil
}

func PrettyPrint(body []byte) error {

	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		Error.Fatalln("error parsing response: ", err)
		return err
	}
	fmt.Println(prettyJSON.String())
	return nil
}

func getPrivateKey() (interface{}, error) {
	pemPrivateKey := fmt.Sprintf("%v", viper.Get("private_key"))
	block, _ := pem.Decode([]byte(pemPrivateKey))
	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		Error.Fatalln("error parsing Private Key: ", err)
		return nil, err
	}
	return privKey, nil
}

func generateJWT() (string, error) {

	const aud = "https://www.googleapis.com/oauth2/v4/token"
	const scope = "https://www.googleapis.com/auth/cloud-platform"

	privKey, err := getPrivateKey()

	if err != nil {
		return "", err
	}

	now := time.Now()
	token := jwt.New()

	_ = token.Set(jwt.AudienceKey, aud)
	_ = token.Set(jwt.IssuerKey, viper.Get("client_email"))
	_ = token.Set("scope", scope)
	_ = token.Set(jwt.IssuedAtKey, now.Unix())
	_ = token.Set(jwt.ExpirationKey, now.Unix())

	payload, err := token.Sign(jwa.RS256, privKey)
	if err != nil {
		Error.Fatalln("error parsing Private Key: ", err)
		return "", err
	}
	Info.Println("jwt token : ", string(payload))
	return string(payload), nil
}

func GenerateAccessToken() (string, error) {

	const tokenEndpoint = "https://www.googleapis.com/oauth2/v4/token"
	const grantType = "urn:ietf:params:oauth:grant-type:jwt-bearer"

	token, err := generateJWT()

	if err != nil {
		return "", nil
	}

	form := url.Values{}
	form.Add("grant_type", grantType)
	form.Add("assertion", token)

	client := &http.Client{}
	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		Error.Fatalln("error in client: ", err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	resp, err := client.Do(req)

	if err != nil {
		Error.Fatalln("failed to generate oauth token: ", err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		Error.Fatalln("error in response: ", string(bodyBytes))
		return "", errors.New("Error in response")
	}
	decoder := json.NewDecoder(resp.Body)
	accessToken := oAuthAccessToken{}
	if err := decoder.Decode(&accessToken); err != nil {
		Error.Fatalln("error in response: ", err)
		return "", errors.New("Error in response")
	}
	Info.Println("access token : ", accessToken)
	RootArgs.Token = accessToken.AccessToken
	_ = writeAccessToken()
	return accessToken.AccessToken, nil
}

func readAccessToken() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(path.Join(usr.HomeDir, accessTokenFile))
	if err != nil {
		Info.Println("Cached access token was not found")
		return err
	}
	Info.Println("Using cached access token: ", string(content))
	RootArgs.Token = string(content)
	return nil
}

func writeAccessToken() error {

	if skipCache {
		return nil
	}

	usr, err := user.Current()
	if err != nil {
		Warning.Println(err)
		return err
	}
	Info.Println("Cache access token: ", RootArgs.Token)
	//don't append access token
	return WriteByteArrayToFile(path.Join(usr.HomeDir, accessTokenFile), false, []byte(RootArgs.Token))
}

func checkAccessToken() bool {

	if skipCheck {
		Info.Println("skipping token validity")
		return true
	}

	const tokenInfo = "https://www.googleapis.com/oauth2/v1/tokeninfo"
	u, _ := url.Parse(tokenInfo)
	q := u.Query()
	q.Set("access_token", RootArgs.Token)
	u.RawQuery = q.Encode()

	client := &http.Client{}

	Info.Println("Connecting to : ", u.String())
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		Error.Fatalln("error in client:", err)
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		Error.Fatalln("error connecting to token endpoint: ", err)
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Error.Fatalln("token info error: ", err)
		return false
	} else if resp.StatusCode != 200 {
		Error.Fatalln("token expired: ", string(body))
		return false
	}
	Info.Println("Response: ", string(body))
	Info.Println("Reusing the cached token: ", RootArgs.Token)
	return true
}

func SetAccessToken() error {

	if RootArgs.Token == "" && RootArgs.ServiceAccount == "" {
		err := readAccessToken() //try to read from config
		if err != nil {
			return fmt.Errorf("either token or service account must be provided")
		}
		if checkAccessToken() { //check if the token is still valid
			return nil
		}
		return fmt.Errorf("token expired: request a new access token or pass the service account")
	}
	if RootArgs.ServiceAccount != "" {
		viper.SetConfigFile(RootArgs.ServiceAccount)
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			return fmt.Errorf("error reading config file: %s", err)
		}
		if viper.Get("private_key") == "" {
			return fmt.Errorf("private key missing in the service account")
		}
		if viper.Get("client_email") == "" {
			return fmt.Errorf("client email missing in the service account")
		}
		_, err = GenerateAccessToken()
		if err != nil {
			return fmt.Errorf("Fatal error generating access token: %s", err)
		}
		return nil
	}
	//a token was passed, cache it
	if checkAccessToken() {
		_ = writeAccessToken()
		return nil
	}
	return fmt.Errorf("token expired: request a new access token or pass the service account")
}

func ReadBundle(filename string) error {

	if !strings.HasSuffix(filename, ".zip") {
		Error.Fatalln("Proxy bundle must be a zip file")
		return errors.New("source must be a zipfile")
	}

	file, err := os.Open(filename)

	if err != nil {
		Error.Fatalln("Cannot open/read API Proxy Bundle: ", err)
		return err
	}

	fi, err := file.Stat()
	if err != nil {
		Error.Fatalln("Error accessing file: ", err)
		return err
	}
	_, err = zip.NewReader(file, fi.Size())

	if err != nil {
		Error.Fatalln("Invalid API Proxy Bundle: ", err)
		return err
	}

	defer file.Close()
	return nil
}

//WriteByteArrayToFile accepts []bytes and writes to a file
func WriteByteArrayToFile(exportFile string, fileAppend bool, payload []byte) error {

	var fileFlags = os.O_CREATE | os.O_WRONLY

	if fileAppend {
		fileFlags = fileFlags | os.O_APPEND
	}

	f, err := os.OpenFile(exportFile, fileFlags, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	//if payload is nil, use internal variable
	if payload != nil {
		_, err = f.Write(payload)
		if err != nil {
			Error.Fatalln("error writing to file: ", err)
			return err
		}
		return nil
	}

	//begin json array
	_, err = f.Write([]byte("["))
	if err != nil {
		Error.Fatalln("error writing to file ", err)
		return err
	}

	payloadFromArray := bytes.Join(EntityPayloadList, []byte(","))
	//add json array terminate
	payloadFromArray = append(payloadFromArray, byte(']'))

	_, err = f.Write(payloadFromArray)

	if err != nil {
		Error.Fatalln("error writing to file: ", err)
		return err
	}

	return nil
}

//GetAsync stores results for each entity in a list
func GetAsyncEntity(entityName string, entityType string, wg *sync.WaitGroup, mu *sync.Mutex) {

	//this is a two step process - 1) get entity details 2) store in byte[][]
	defer wg.Done()
	var respBody []byte

	u, _ := url.Parse(BaseURL)
	u.Path = path.Join(u.Path, RootArgs.Org, entityType, entityName)
	//don't print to sysout
	respBody, err := HttpClient(false, u.String())

	if err != nil {
		Error.Fatalln("Error with entity: %s", entityName)
		Error.Fatalln(err)
		return
	}

	mu.Lock()
	EntityPayloadList = append(EntityPayloadList, respBody)
	mu.Unlock()
	Info.Printf("Completed entity: %s", entityName)
}

//FetchAyncBundle can download a shared flow or a proxy bundle
func FetchAsyncBundle(entityType string, name string, revision string, wg *sync.WaitGroup) {

	defer wg.Done()

	u, _ := url.Parse(BaseURL)
	q := u.Query()
	q.Set("format", "bundle")
	u.RawQuery = q.Encode()
	u.Path = path.Join(u.Path, RootArgs.Org, entityType, name, "revisions", revision)

	err := DownloadResource(u.String(), name)
	if err != nil {
		Error.Fatalln("Error with entity: %s", name)
		Error.Fatalln(err)
	}
}
