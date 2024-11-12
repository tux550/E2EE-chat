package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	e2ee_api "tux.tech/e2ee/api"
	x3dh_client "tux.tech/x3dh/client"
	x3dh_core "tux.tech/x3dh/core"
)

// ================================== CONFIG ===========================
var url = "wss://localhost:8765/ws"
var contacts_filename = "contacts.json"
var secrets_filename string

// ================================== PRETTY PRINT ===========================
func prettyAskString(question string) string {
	fmt.Print(text.FgGreen.Sprintf(question))
	var answer string
	fmt.Scanln(&answer)
	return answer
}

func prettyAskInt(question string) int {
	fmt.Print(text.FgGreen.Sprintf(question))
	var answer int
	fmt.Scanln(&answer)
	return answer
}

func prettyLogInfo(info string) {
	fmt.Println(text.FgHiBlack.Sprintf(info))
}

func prettyLogRisky(info string) {
	fmt.Println(text.FgHiRed.Sprintf(info))
}

func prettyTitle(title string) {
	fmt.Println(text.FgHiCyan.Sprintf(title))
}

// ================================== STRUCTS ===========================
type Contact struct {
	Username  string
	PublicKey x3dh_core.X3DHPublicIK
}

func (c Contact) PrettyPrint() {
	// Create and configure the table writer
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Field", "Value"})

	// Add rows to the table
	t.AppendRows([]table.Row{
		{"Username", c.Username},
		{"Public Key", base64.StdEncoding.EncodeToString(c.PublicKey.IdentityKey[:])},
	})

	// Customize table appearance
	t.SetStyle(table.StyleColoredBright)
	t.Style().Format.Header = text.FormatDefault
	t.Style().Options.SeparateRows = true

	// Print the title and render the table
	prettyTitle("=== Contact ===")
	t.Render()
}

func (c Contact) DebugPrint() {
	fmt.Println("=== Contact ===")
	fmt.Println("Username:", c.Username)
	fmt.Println("Public Key:", base64.StdEncoding.EncodeToString(c.PublicKey.IdentityKey[:]))
	fmt.Println("===============")
}

func (c Contact) ExportToFile(filename string) error {
	// Marshal the contact to JSON
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	// Write the data to the file
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	// Return
	return nil
}

func ImportContactFromFile(filename string) (*Contact, error) {
	// Read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// Unmarshal the data
	var contact Contact
	err = json.Unmarshal(data, &contact)
	if err != nil {
		return nil, err
	}
	// Return
	return &contact, nil
}

type Contacts []Contact

func (c *Contacts) AddContact(contact Contact) error {
	// If public key is already in contacts, return error
	for _, c := range *c {
		if c.PublicKey.IdentityKey.Equal(contact.PublicKey.IdentityKey) {
			return fmt.Errorf("Contact already exists")
		}
	}
	*c = append(*c, contact)
	return nil
}

func (c *Contacts) RemoveContact(id int) {
	*c = append((*c)[:id], (*c)[id+1:]...)
}

func (c Contacts) GetContact(id int) Contact {
	return c[id]
}

func (c Contacts) FindContactByUsername(username string) *Contact {
	for _, contact := range c {
		if contact.Username == username {
			return &contact
		}
	}
	return nil
}

func InitContacts() *Contacts {
	return &Contacts{}
}

func SaveContacts(contacts *Contacts, filename string) error {
	// Marshal the contacts to JSON
	data, err := json.Marshal(contacts)
	if err != nil {
		return err
	}
	// Write the data to the file
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	// Return
	return nil
}

func LoadContacts(filename string) (*Contacts, error) {
	// Read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// Unmarshal the data
	var contacts Contacts
	err = json.Unmarshal(data, &contacts)
	if err != nil {
		return nil, err
	}
	// Return
	return &contacts, nil
}

// ================================== HELPER ===========================
func GetMyContact(client *x3dh_client.X3DHClient) Contact {
	return Contact{
		Username:  client.Username,
		PublicKey: *client.IdentityKey.PublicIK(),
	}
}

// ================================== CLIENT MANAGER ===========================
func GetMyClient() (*x3dh_client.X3DHClient, error) {
	// Check if secrets file exists
	_, err := os.Stat(secrets_filename)
	if os.IsNotExist(err) {
		// Create new secrets client
		prettyLogInfo("No existing client found")
		prettyLogInfo("Creating new client")
		username := prettyAskString("Enter username: ")

		client, err := x3dh_client.InitClient(username)
		if err != nil {
			return nil, err
		}
		// Save secrets to file
		err = client.SaveClient(secrets_filename)
		if err != nil {
			return nil, err
		}
		return client, nil
	} else {
		// Load secrets from file
		prettyLogInfo("Loaded existing client")
		client, err := x3dh_client.LoadClient(secrets_filename)
		if err != nil {
			return nil, err
		}
		return client, nil
	}
}

func GetMyContacts() (*Contacts, error) {
	// Check if contacts file exists
	_, err := os.Stat(contacts_filename)
	if os.IsNotExist(err) {
		// Create new contacts
		prettyLogInfo("No existing contacts found")
		prettyLogInfo("Creating new contacts")
		contacts := InitContacts()
		// Save contacts to file
		err = SaveContacts(contacts, contacts_filename)
		if err != nil {
			return nil, err
		}
		return contacts, nil
	} else {
		// Load contacts from file
		prettyLogInfo("Loaded existing contacts")
		contacts, err := LoadContacts(contacts_filename)
		if err != nil {
			return nil, err
		}
		return contacts, nil
	}
}

func SaveMyClient(client *x3dh_client.X3DHClient) error {
	// Save secrets to file
	//fmt.Println("Saved client")
	err := client.SaveClient(secrets_filename)
	if err != nil {
		return err
	}
	return nil
}

func SaveMyContacts(contacts *Contacts) error {
	// Save contacts to file
	//fmt.Println("Saved contacts")
	err := SaveContacts(contacts, contacts_filename)
	if err != nil {
		return err
	}
	return nil
}

// ================================== HANDLE INCOMING MESSAGES ===========================
/*
func HandleIncomingMessages(client *x3dh_client.X3DHClient, c *websocket.Conn) {
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println("Could not read message:", err)
			return
		}
		if mt == websocket.TextMessage {
			fmt.Println("Received message:", string(message))
		}
	}
}
*/

// ================================== Options ===========================
func MenuListContacts(contacts *Contacts) {
	// Add title to table
	fmt.Println(text.FgHiCyan.Sprintf("=== Contacts ==="))
	if len(*contacts) == 0 {
		fmt.Println("No contacts")
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Username", "Public Key"})

	for i, contact := range *contacts {
		t.AppendRow([]interface{}{
			i,
			contact.Username,
			base64.StdEncoding.EncodeToString(contact.PublicKey.IdentityKey[:]),
		})
	}

	t.SetStyle(table.StyleColoredBright)
	t.Style().Format.Header = text.FormatDefault
	t.Style().Options.SeparateRows = true

	t.Render()

}

func MenuAddContact(client *x3dh_client.X3DHClient, contacts *Contacts) {
	// Read contact from file
	filename := prettyAskString("Enter contact file: ")

	contact, err := ImportContactFromFile(filename)
	if err != nil {
		prettyLogRisky("Could not import contact from file")
		//fmt.Println("Could not import contact from file:", err)
		return
	}
	// Add contact to contacts
	err = contacts.AddContact(*contact)
	if err != nil {
		prettyLogRisky("Contact already exists")
		//fmt.Println("Contact already exists:", err)
		return
	}
	// Save contacts
	err = SaveMyContacts(contacts)
	if err != nil {
		prettyLogRisky("Could not save contacts")
		//fmt.Println("Could not save contacts:", err)
		return
	}
	prettyLogInfo("Contact added")
	//fmt.Println("Contact added")
}

func MenuRemoveContact(contacts *Contacts) {
	// Read contact id
	fmt.Println("Enter contact id:")
	var id int
	fmt.Scanln(&id)
	// Confirmation
	fmt.Println("Are you sure you want to remove contact:")
	contacts.GetContact(id).DebugPrint()
	fmt.Println("Enter 'yes' to confirm:")
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "yes" {
		fmt.Println("Contact not removed")
		return
	}
	// Remove contact
	contacts.RemoveContact(id)
	// Save contacts
	err := SaveMyContacts(contacts)
	if err != nil {
		fmt.Println("Could not save contacts:", err)
		return
	}
	fmt.Println("Contact removed")
}

func MenuChat(client *x3dh_client.X3DHClient, contacts *Contacts, c *websocket.Conn) {
	// Select contact
	id := prettyAskInt("Enter contact id: ")

	contact := contacts.GetContact(id)
	// Write message
	message := prettyAskString("Enter message: ")
	// Send message
	success, err := APISendMessage(client, c, contact, []byte(message))
	if err != nil {
		prettyLogRisky("Could not send message")
		return
	}
	if !success {
		prettyLogRisky("Could not send message")
		return
	}
	// Success
	prettyLogInfo("Message sent")
}

func MenuShareMyContact(client *x3dh_client.X3DHClient) {
	// Get my contact
	contact := GetMyContact(client)
	// Pretty print my contact
	contact.PrettyPrint()
	// Export my contact to file
	err := contact.ExportToFile("MyContact.json")
	if err != nil {
		prettyLogRisky("Could not export contact to file")
		return
	}
	prettyLogInfo("Contact exported to MyContact.json")
}

func MenuReceiveMessages(client *x3dh_client.X3DHClient, c *websocket.Conn, contacts *Contacts) {
	for {
		// Receive message
		message, sender, err := APIReceiveMessage(client, c)
		if err != nil {
			prettyLogRisky("Could not receive message")
			//fmt.Println("Could not receive message:", err)
			return
		}
		if message == nil {
			prettyLogInfo("No more messages")
			return
		}
		// Get contact
		contact := contacts.FindContactByUsername(sender)
		if contact == nil {
			prettyLogRisky("Be cautious, the following message is from an unknown contact: " + sender)
			//fmt.Println("The following message is from an unknown contact: ", sender)
		}
		// Decrypt message
		plaintext, err := client.RecieveMessage(message)
		if err != nil {
			prettyLogRisky("Failed to decrypt message from: " + sender)
			// Continue to next message
			continue
		}
		// Print message
		/*
			fmt.Println("=== Message ===")
			fmt.Println("Sender:", sender)
			fmt.Println("Message:", string(plaintext))
			fmt.Println("===============")
			fmt.Println()*/
		// Create and configure the table writer
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)

		// Add rows to the table
		t.AppendRows([]table.Row{
			{"Sender", sender},
			{"Message", string(plaintext)},
		})

		// Customize table appearance
		t.SetStyle(table.StyleColoredBright)
		t.Style().Format.Header = text.FormatDefault
		t.Style().Options.SeparateRows = true
		t.Style().Options.SeparateColumns = false

		// Print the title and render the table
		prettyTitle("=== Message ===")
		t.Render()
		fmt.Println()
	}
}

func MenuHelp() {
	fmt.Println()
	fmt.Printf("=== Welcome to the E2EE Client ===\n")
	fmt.Printf("This is a simple CLI client for end-to-end encryption.\n")
	fmt.Printf("It uses the X3DH protocol to establish a secure connection.\n")

	fmt.Println()
	fmt.Println("=== Security Notice ===")
	fmt.Println("This client ONLY guarantees secure communication between clients that have previously exchanged and verified their public keys.")
	fmt.Println("It is the user's responsibility to ensure that the public keys are correct and have not been tampered with.")

	fmt.Println()
	fmt.Printf("=== Menu Options ===\n")
	fmt.Println("List Contacts: List all contacts")
	fmt.Println("Add Contact: Add a new contact")
	fmt.Println("Remove Contact: Remove a contact")
	fmt.Println("Send Message: Send a message to a contact")
	fmt.Println("Receive Messages: Receive all messages")
	fmt.Println("Share My Contact: Export my contact to a file")
	fmt.Println("Exit: Exit the program")
}

// ================================== User Interface ===========================
func showMenu() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Options "})

	menuItems := []struct {
		Index int
		Item  string
	}{
		{1, "List Contacts"},
		{2, "Add Contact"},
		{3, "Remove Contact"},
		{4, "Send Message"},
		{5, "Receive Messages"},
		{6, "Share My Contact"},
		{7, "Help"},
		{8, "Exit"},
	}

	for _, menuItem := range menuItems {
		t.AppendRow([]interface{}{menuItem.Index, menuItem.Item})
	}

	t.SetStyle(table.StyleColoredBright)

	// Add padding before render
	fmt.Println()
	fmt.Println(text.FgHiCyan.Sprintf("=== Menu ==="))
	t.Render()
}

func Menu(client *x3dh_client.X3DHClient, contacts *Contacts, c *websocket.Conn) {
	for {
		showMenu()
		choice := prettyAskInt("Enter choice: ")
		// Add padding after choice
		fmt.Println()

		switch choice {
		case 1:
			MenuListContacts(contacts)
		case 2:
			MenuAddContact(client, contacts)
		case 3:
			MenuRemoveContact(contacts)
		case 4:
			MenuChat(client, contacts, c)
		case 5:
			MenuReceiveMessages(client, c, contacts)
		case 6:
			MenuShareMyContact(client)
		case 7:
			MenuHelp()
		case 8:
			fmt.Println("Exit")
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

// ================================== API CALLS ===========================
func sendAndAwaitWsResponse(c *websocket.Conn, params interface{}, method string) (json.RawMessage, error) {
	// Marshal params
	marshalledParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	// Build API call
	api_call := &e2ee_api.InboundMessage{
		Method: method,
		Params: marshalledParams,
	}
	// Marshal
	data, err := json.Marshal(api_call)
	if err != nil {
		return nil, err
	}
	// Send request
	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return nil, err
	}
	// Await response from server
	// (message is in channel incomingResponses)
	response := <-incomingResponses

	// Check for success
	if response.Method != method {
		fmt.Println("Received wrong method. Expected:", method, "Received:", response.Method)
		return nil, fmt.Errorf("received wrong method")
	}

	// DEBUG: Print recieved message
	// fmt.Println("Received message:", string(message))

	// Return params
	return response.Params, nil
}

func APIUploadBundle(client *x3dh_client.X3DHClient, c *websocket.Conn) (bool, error) {
	// Get bundle
	bundle, err := client.GetServerInitBundle()
	if err != nil {
		return false, err
	}
	// Build API call
	params := &e2ee_api.RequestUploadBundle{
		UserID: client.Username,
		Bundle: *bundle,
	}
	// Send And Await Response
	response, err := sendAndAwaitWsResponse(c, params, "upload_bundle")
	if err != nil {
		return false, err
	}
	// Parse params
	params_response := &e2ee_api.ResponseUploadBundle{}
	err = json.Unmarshal(response, params_response)
	if err != nil {
		return false, err
	}
	// Return status
	return params_response.Success, nil
}

func APIGetBundle(client *x3dh_client.X3DHClient, c *websocket.Conn, contact Contact) (*x3dh_core.X3DHKeyBundle, error) {
	// Build API call
	params := &e2ee_api.RequestUserBundle{
		UserID: contact.Username,
	}
	// Send And Await Response
	response, err := sendAndAwaitWsResponse(c, params, "get_bundle")
	if err != nil {
		return nil, err
	}
	// Parse params
	params_response := &e2ee_api.ResponseUserBundle{}
	err = json.Unmarshal(response, params_response)
	if err != nil {
		return nil, err
	}
	// Success
	if !params_response.Success {
		return nil, fmt.Errorf("failed to get bundle")
	}
	// Validate bundle
	if !params_response.Bundle.IK.IdentityKey.Equal(contact.PublicKey.IdentityKey) {
		return nil, fmt.Errorf("bundle identity key does not match contact public key")
	}
	if !params_response.Bundle.Validate() {
		return nil, fmt.Errorf("failed to validate bundle")
	}
	// Return status
	return &params_response.Bundle, nil
}

func APISendMessage(client *x3dh_client.X3DHClient, c *websocket.Conn, contact Contact, message []byte) (bool, error) {
	// Get contact bundle
	bundle, err := APIGetBundle(client, c, contact)
	if err != nil {
		return false, err
	}
	// Encrypt message
	x3dhMessage, err := client.BuildMessage(bundle, message)
	if err != nil {
		return false, err
	}
	// Build API call
	params := &e2ee_api.RequestSendMsg{
		RecipientID: contact.Username,
		MessageData: *x3dhMessage,
	}
	// Send And Await Response
	response, err := sendAndAwaitWsResponse(c, params, "send_message")
	if err != nil {
		return false, err
	}
	// Parse params
	params_response := &e2ee_api.ResponseSendMsg{}
	err = json.Unmarshal(response, params_response)
	if err != nil {
		return false, err
	}
	// Return status
	return params_response.Success, nil
}

func APIReceiveMessage(client *x3dh_client.X3DHClient, c *websocket.Conn) (*x3dh_core.InitialMessage, string, error) {
	// Send And Await Response
	response, err := sendAndAwaitWsResponse(c, e2ee_api.RequestReceiveMsg{}, "receive_message")
	if err != nil {
		return nil, "", err
	}
	// Parse params
	params_response := &e2ee_api.ResponseReceiveMsg{}
	err = json.Unmarshal(response, params_response)
	if err != nil {
		return nil, "", err
	}
	// End of queue
	if !params_response.Success {
		return nil, "", nil
	}
	// Return message data
	return &params_response.MessageData, params_response.SenderID, nil
}

func APIGetStatus(client *x3dh_client.X3DHClient, c *websocket.Conn) (bool, error) {
	// Build API call
	params := &e2ee_api.RequestUserStatus{}
	// Send And Await Response
	response, err := sendAndAwaitWsResponse(c, params, "status")
	if err != nil {
		return false, err
	}
	// Parse params
	params_response := &e2ee_api.ResponseUserStatus{}
	err = json.Unmarshal(response, params_response)
	if err != nil {
		return false, err
	}
	// Return status
	return params_response.Success, nil
}

// ================================== CONNECTION ===========================
// Setup channels for incoming messages
var incomingResponses = make(chan e2ee_api.OutboundMessage)
var incomingNotifications = make(chan e2ee_api.OutboundMessage)

func ReadIncomingMessages(c *websocket.Conn) {
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			return
		}
		if mt == websocket.TextMessage {
			// Unmarshal the message
			response := &e2ee_api.OutboundMessage{}
			err = json.Unmarshal(message, response)
			if err != nil {
				prettyLogRisky("Recieved invalid data from server")
				continue
			}
			// Check if message is a notification (Method starts with notify_)
			if strings.HasPrefix(response.Method, "notify_") {
				incomingNotifications <- *response
			} else {
				incomingResponses <- *response
			}
		}
	}
}

func APIUploadNewOTPs(client *x3dh_client.X3DHClient, c *websocket.Conn) (bool, error) {
	// Get OTPs
	otps, err := client.BatchGenerateOTPs(5)
	// Save client after generating OTPs
	SaveMyClient(client)
	if err != nil {
		return false, err
	}

	// Build API call
	params := &e2ee_api.RequestUploadOTPs{
		OTPs: otps,
	}

	// Marshal params
	marshalledParams, err := json.Marshal(params)
	if err != nil {
		return false, err
	}

	// Build API call
	api_call := &e2ee_api.InboundMessage{
		Method: "upload_new_otps",
		Params: marshalledParams,
	}

	// Marshal
	data, err := json.Marshal(api_call)
	if err != nil {
		return false, err
	}

	// Send request
	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return false, err
	}

	return true, nil
}

func HandleNotifications(client *x3dh_client.X3DHClient, c *websocket.Conn) {
	for {
		// Wait for notification
		notification := <-incomingNotifications
		// Handle notification
		switch notification.Method {
		case "notify_low_otp":
			fmt.Println()
			prettyLogInfo("Low OTP. Sending more.")
			APIUploadNewOTPs(client, c)
			fmt.Println()
		case "notify_new_message":
			fmt.Println()
			prettyLogInfo("<New message pending>")
			fmt.Println()
		default:
			fmt.Println()
			prettyLogRisky("Unknown notification")
			fmt.Println()
		}
	}
}

func ConnectToServer(client *x3dh_client.X3DHClient, contacts *Contacts) {
	// Set username as header
	header := http.Header{}
	header.Add("User", client.Username)

	// Get password
	password := prettyAskString("Enter password: ")
	header.Add("Password", password)

	// Load the server's certificate
	caCert, err := ioutil.ReadFile("../certs/server.crt")
	if err != nil {
		fmt.Println("Error reading CA certificate:", err)
		return
	}

	// Create a certificate pool and add the server's certificate
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a custom dialer with TLS configuration if needed
	dialer := websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{
		//RootCAs: caCertPool,
		InsecureSkipVerify: true, // Set to true if you're using a self-signed certificate for testing
	}

	// Connect to server
	c, _, err := dialer.Dial(url, header)
	if err != nil {
		//prettyLogRisky("Connection failed")
		fmt.Println("Could not connect to server:", err)
		return
	}

	// Start reading incoming messages
	go ReadIncomingMessages(c)

	// Deferred close
	defer func() {
		// Close connection
		c.Close()
		fmt.Println("Closed connection")
		// Save client status
		err := SaveMyClient(client)
		if err != nil {
			fmt.Println("FATAL ERROR - Could not save client:", err)
			return
		}
		//fmt.Println("Saved client and contacts")
	}()

	// Handle notifications in background
	go HandleNotifications(client, c)

	// Check Status
	success, err := APIGetStatus(client, c)
	if err != nil {
		prettyLogRisky("Could not get status")
		fmt.Println("Could not get status:", err)
		return
	}

	if !success {
		prettyLogInfo("First time setup")
		success, err := APIUploadBundle(client, c)
		if err != nil {
			prettyLogRisky("Could not upload bundle")
			fmt.Println("Could not upload bundle:", err)
			return
		}
		if !success {
			prettyLogRisky("Could not upload bundle")
			fmt.Println("Could not upload bundle 2")
			return
		}
	}
	// Infinite loop for interface
	Menu(client, contacts, c)
}

// ================================== MAIN ===========================
func main() {
	// Select user file
	secrets_filename = prettyAskString("Enter user file: ")

	// LOAD CLIENT
	a, err := GetMyClient()
	if err != nil {
		prettyLogRisky("Failed to load client!")
		return
	}
	// LOAD CONTACTS
	c, err := GetMyContacts()
	if err != nil {
		prettyLogRisky("Failed to load contacts!")
		return
	}
	// Connect to server
	ConnectToServer(a, c)
}
