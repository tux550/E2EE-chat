package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	e2ee_api "tux.tech/e2ee/api"
	x3dh_client "tux.tech/x3dh/client"
	x3dh_core "tux.tech/x3dh/core"
)

// ================================== CONFIG ===========================
var url = "ws://localhost:8765/ws"
var contacts_filename = "contacts.json"
var secrets_filename string

// ================================== STRUCTS ===========================
type Contact struct {
	Username  string
	PublicKey x3dh_core.X3DHPublicIK
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
		fmt.Println("Initialized new client")
		fmt.Println("Enter your username: ")
		var username string
		fmt.Scanln(&username)
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
		fmt.Println("Loaded existing client")
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
		fmt.Println("Initialized new contacts")
		contacts := InitContacts()
		// Save contacts to file
		err = SaveContacts(contacts, contacts_filename)
		if err != nil {
			return nil, err
		}
		return contacts, nil
	} else {
		// Load contacts from file
		fmt.Println("Loaded existing contacts")
		contacts, err := LoadContacts(contacts_filename)
		if err != nil {
			return nil, err
		}
		return contacts, nil
	}
}

func SaveMyClient(client *x3dh_client.X3DHClient) error {
	// Save secrets to file
	fmt.Println("Saved client")
	err := client.SaveClient(secrets_filename)
	if err != nil {
		return err
	}
	return nil
}

func SaveMyContacts(contacts *Contacts) error {
	// Save contacts to file
	fmt.Println("Saved contacts")
	err := SaveContacts(contacts, contacts_filename)
	if err != nil {
		return err
	}
	return nil
}

// ================================== HANDLE INCOMING MESSAGES ===========================
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

// ================================== Options ===========================
func MenuListContacts(contacts *Contacts) {
	fmt.Println("=== Contacts ===")
	for i, contact := range *contacts {
		fmt.Println(i, contact.Username, base64.StdEncoding.EncodeToString(contact.PublicKey.IdentityKey[:]))
	}
	fmt.Println("===============")
}

func MenuAddContact(client *x3dh_client.X3DHClient, contacts *Contacts) {
	// Read contact from file
	fmt.Println("Enter contact file:")
	var filename string
	fmt.Scanln(&filename)
	contact, err := ImportContactFromFile(filename)
	if err != nil {
		fmt.Println("Could not import contact from file:", err)
		return
	}
	// Add contact to contacts
	err = contacts.AddContact(*contact)
	if err != nil {
		fmt.Println("Contact already exists:", err)
		return
	}
	// Save contacts
	err = SaveMyContacts(contacts)
	if err != nil {
		fmt.Println("Could not save contacts:", err)
		return
	}
	fmt.Println("Contact added")
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

func MenuSendMessage(client *x3dh_client.X3DHClient, contacts *Contacts) {
	// Read contact id
	fmt.Println("Enter contact id:")
	var id int
	fmt.Scanln(&id)
	// Get contact
	contact := contacts.GetContact(id)
	fmt.Println("Sending message to:")
	contact.DebugPrint()
	// Read message
	fmt.Println("Enter message:")
	var message string
	fmt.Scanln(&message)
	// TODO
	// Get contact data from server
	// Encrypt message
	// Send message
}

func MenuShareMyContact(client *x3dh_client.X3DHClient) {
	// Get my contact
	contact := GetMyContact(client)
	// Export my contact to file
	err := contact.ExportToFile("MyContact.json")
	if err != nil {
		fmt.Println("Could not export contact to file:", err)
		return
	}
	fmt.Println("Contact exported to MyContact.json")
}

// ================================== User Interface ===========================
func Menu(client *x3dh_client.X3DHClient, contacts *Contacts, c *websocket.Conn) {
	for {
		fmt.Println("=== Menu ===")
		fmt.Println("1. List Contacts")
		fmt.Println("2. Add Contact")
		fmt.Println("3. Remove Contact")
		fmt.Println("4. Send Message")
		fmt.Println("5. Share My Contact")
		fmt.Println("6. Exit")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Println("List Contacts")
			MenuListContacts(contacts)
		case 2:
			fmt.Println("Add Contact")
			MenuAddContact(client, contacts)
		case 3:
			fmt.Println("Remove Contact")
			MenuRemoveContact(contacts)
		case 4:
			fmt.Println("Send Message")
		case 5:
			fmt.Println("Share My Contact")
			MenuShareMyContact(client)
		case 6:
			fmt.Println("Exit")
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

// ================================== API CALLS ===========================

func APIUploadBundle(client *x3dh_client.X3DHClient, c *websocket.Conn) {
	// Get bundle
	bundle, err := client.GetServerInitBundle()
	if err != nil {
		fmt.Println("Could not get bundle:", err)
		return
	}
	// Build API call
	params_data := &e2ee_api.RequestUploadBundle{
		UserID: client.Username,
		Bundle: *bundle,
	}
	params, err := json.Marshal(params_data)
	if err != nil {
		fmt.Println("Could not marshal bundle:", err)
		return
	}
	api_call := &e2ee_api.InboundMessage{
		Method: "upload_bundle",
		Params: params,
	}
	// Marshal the API call to JSON
	data, err := json.Marshal(api_call)
	if err != nil {
		fmt.Println("Could not marshal API call:", err)
		return
	}
	// Send bundle
	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		fmt.Println("Could not send bundle:", err)
		return
	}
}

// ================================== CONNECTION ===========================

func ConnectToServer(client *x3dh_client.X3DHClient, contacts *Contacts) {
	// Set username as header
	header := http.Header{}
	header.Add("User", client.Username)

	// Connect to server
	c, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		fmt.Println("Could not connect to server:", err)
		return
	}

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

	// Parallel handle incoming messages
	//go HandleIncomingMessages(client, c)

	// Infinite loop for interface
	APIUploadBundle(client, c)
	Menu(client, contacts, c)
	/*
			// Read user input
			var message string
			fmt.Scanln(&message)

			// Check for exit
			if message == "exit" {
				return
			}

			// Send message
			err = c.WriteMessage(websocket.BinaryMessage, []byte(message))
			if err != nil {
				fmt.Println("Could not send message:", err)
				return
			}
			fmt.Println("Sent message:", message)
		}
	*/
}

// ================================== MAIN ===========================
func main() {
	// Select user file
	fmt.Println("Enter user file:")
	fmt.Scanln(&secrets_filename)

	// LOAD CLIENT
	a, err := GetMyClient()
	if err != nil {
		fmt.Println("Could not get client:", err)
		return
	}
	// LOAD CONTACTS
	c, err := GetMyContacts()
	if err != nil {
		fmt.Println("Could not get contacts:", err)
		return
	}

	// Debug Print Client
	a.DebugPrint()

	// Connect to server
	ConnectToServer(a, c)
}
