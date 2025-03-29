
package client_test

// You MUST NOT change these default imports.  ANY additional imports may
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports.
	_ "encoding/hex"
	_ "errors"
	_ "strconv"
	_ "strings"
	"testing"

	_ "github.com/google/uuid"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect().
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"

	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Tests")
}

// ================================================
// Global Variables (feel free to add more!)
// ================================================
const defaultPassword = "password"
const defaultPassword2 = "password2"
const emptyString = ""
const contentOne = "Bitcoin is Nick's favorite "
const contentTwo = "digital "
const contentThree = "cryptocurrency!"

// ================================================
// Describe(...) blocks help you organize your tests
// into functional categories. They can be nested into
// a tree-like structure.
// ================================================

var _ = Describe("Client Tests", func() {

	// A few user declarations that may be used for testing. Remember to initialize these before you
	// attempt to use them!
	var alice *client.User
	var bob *client.User
	var charles *client.User
	// var doris *client.User
	// var eve *client.User
	// var frank *client.User
	// var grace *client.User
	// var horace *client.User
	// var ira *client.User

	// These declarations may be useful for multi-session testing.
	var alicePhone *client.User
	var aliceLaptop *client.User
	var aliceDesktop *client.User

	var err error

	// A bunch of filenames that may be useful.
	aliceFile := "aliceFile.txt"
	bobFile := "bobFile.txt"
	charlesFile := "charlesFile.txt"

	// Newly Added!
	fileName := "sharedfile.txt"
	initialContent := "file from laptop"
	phoneAppend := " + updated from phone"
	aliceUpdate := "Update by Alice. "
	bobUpdate := "Update by Bob."
	bobFileAlias := "bob_view.txt"
	content := "mission-critical data"

	// dorisFile := "dorisFile.txt"
	// eveFile := "eveFile.txt"
	// frankFile := "frankFile.txt"
	// graceFile := "graceFile.txt"
	// horaceFile := "horaceFile.txt"
	// iraFile := "iraFile.txt"

	BeforeEach(func() {
		// This runs before each test within this Describe block (including nested tests).
		// Here, we reset the state of Datastore and Keystore so that tests do not interfere with each other.
		// We also initialize
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	Describe("Basic Tests", func() {

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Alice again.")
			alice, err = client.InitUser("alice", defaultPassword2)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Alice trying to login with incorrect password:%s", defaultPassword2)
			_, err = client.GetUser("alice", defaultPassword2)
			Expect(err).ToNot(BeNil())
		})

		Specify("Password Validation: Invalid credentials should return an error", func() {
			userlib.DebugMsg("Initializing user Alice with defaultPassword.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Attempting GetUser with incorrect password.")
			_, err = client.GetUser("alice", "hahahaha")
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Attempting to create a user with the same username.")
			_, err = client.InitUser("alice", "anotherPassword")
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Attempting to log in with a different username")
			_, err = client.GetUser("bob", defaultPassword)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Verifying that correct password still works.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			Expect(aliceLaptop).ToNot(BeNil())
		})

		Specify("Basic Test: Testing Single User Store/Load/Append.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Multiple Devices: Users on different devices see synchronized changes", func() {
			userlib.DebugMsg("Initializing user EvanBot on laptop.")
			evanLaptop, err := client.InitUser("evanbot", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting EvanBot on phone (another device).")
			evanPhone, err := client.GetUser("evanbot", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Laptop stores a new file.")
			err = evanLaptop.StoreFile(fileName, []byte(initialContent))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Phone loads the file and verifies content.")
			content, err := evanPhone.LoadFile(fileName)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(initialContent)))

			userlib.DebugMsg("Phone appends to the file.")
			err = evanPhone.AppendToFile(fileName, []byte(phoneAppend))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Laptop loads file and should see full content.")
			content, err = evanLaptop.LoadFile(fileName)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(initialContent + phoneAppend)))
		})

		Specify("File Sharing: Invitation creation and access via AcceptInvitation", func() {

			userlib.DebugMsg("Initializing users Alice (owner) and Bob (recipient).")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice stores a file.")
			err = alice.StoreFile(fileName, []byte(initialContent))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creates an invitation for Bob.")
			invitationPtr, err := alice.CreateInvitation(fileName, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts the invitation.")
			err = bob.AcceptInvitation("alice", invitationPtr, bobFileAlias)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob loads the shared file and verifies initial content.")
			data, err := bob.LoadFile(bobFileAlias)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(initialContent)))

			userlib.DebugMsg("Alice appends content to the file.")
			err = alice.AppendToFile(fileName, []byte(aliceUpdate))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob reads the file and sees Alice’s update.")
			data, err = bob.LoadFile(bobFileAlias)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(initialContent + aliceUpdate)))

			userlib.DebugMsg("Bob appends content to the file.")
			err = bob.AppendToFile(bobFileAlias, []byte(bobUpdate))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice loads the file and sees all updates.")
			data, err = alice.LoadFile(fileName)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(initialContent + aliceUpdate + bobUpdate)))

		})

		Specify("File Revocation: Only intended users lose access", func() {
			userlib.DebugMsg("Initializing users A (owner), B, C, D, E, F, G.")

			alice, _ := client.InitUser("alice", defaultPassword)
			bob, _ := client.InitUser("bob", defaultPassword)
			charles, _ := client.InitUser("charles", defaultPassword)
			doris, _ := client.InitUser("doris", defaultPassword)
			ella, _ := client.InitUser("ella", defaultPassword)
			grace, _ := client.InitUser("grace", defaultPassword)

			// A creates file
			alice.StoreFile(fileName, []byte(content))

			// A → B
			inviteB, _ := alice.CreateInvitation(fileName, "bob")
			bob.AcceptInvitation("alice", inviteB, "b_file")

			// B → D, E
			inviteD, _ := bob.CreateInvitation("b_file", "doris")
			doris.AcceptInvitation("bob", inviteD, "d_file")

			inviteE, _ := bob.CreateInvitation("b_file", "ella")
			ella.AcceptInvitation("bob", inviteE, "e_file")

			// A → C
			inviteC, _ := alice.CreateInvitation(fileName, "charles")
			charles.AcceptInvitation("alice", inviteC, "c_file")

			// C → G
			inviteG, _ := charles.CreateInvitation("c_file", "grace")
			grace.AcceptInvitation("charles", inviteG, "g_file")

			userlib.DebugMsg("Alice revokes Bob.")
			err := alice.RevokeAccess(fileName, "bob")
			Expect(err).To(BeNil())

			// Confirm that Bob, D, E, and F cannot access
			_, err = bob.LoadFile("b_file")
			Expect(err).ToNot(BeNil())

			_, err = doris.LoadFile("d_file")
			Expect(err).ToNot(BeNil())

			_, err = ella.LoadFile("e_file")
			Expect(err).ToNot(BeNil())

			// Confirm that Charles and Grace still can
			data, err := charles.LoadFile("c_file")
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(content)))

			data, err = grace.LoadFile("g_file")
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(content)))

			// Confirm revoked users cannot append or overwrite
			err = bob.AppendToFile("b_file", []byte("tamper attempt"))
			Expect(err).ToNot(BeNil())

			err = doris.StoreFile("d_file", []byte("overwrite"))
			Expect(err).ToNot(BeNil())
		})

		Specify("AcceptInvitation fails if filename already exists", func() {
			alice, _ = client.InitUser("alice", defaultPassword)
			bob, _ = client.InitUser("bob", defaultPassword)

			alice.StoreFile("report.txt", []byte("midterm"))
			bob.StoreFile("report.txt", []byte("different file"))

			invite, _ := alice.CreateInvitation("report.txt", "bob")
			err = bob.AcceptInvitation("alice", invite, "report.txt")

			Expect(err).ToNot(BeNil()) // Cannot overwrite own file
		})

		Specify("Revoked invitation before acceptance is invalid", func() {
			alice, _ = client.InitUser("alice", defaultPassword)
			bob, _ = client.InitUser("bob", defaultPassword)

			alice.StoreFile("data.txt", []byte("do not share"))
			invite, _ := alice.CreateInvitation("data.txt", "bob")
			alice.RevokeAccess("data.txt", "bob")

			err := bob.AcceptInvitation("alice", invite, "bobdata.txt")
			Expect(err).ToNot(BeNil()) // Should fail
		})

		Specify("Empty filename rejected", func() {
			alice, _ = client.InitUser("alice", defaultPassword)
			err := alice.StoreFile("", []byte("hi"))
			Expect(err).ToNot(BeNil())
		})

		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			aliceDesktop, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob appending to file %s, content: %s", bobFile, contentTwo)
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop appending to file %s, content: %s", aliceFile, contentThree)
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Getting third instance of Alice - alicePhone.")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alicePhone sees Alice's changes.")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Revoke Functionality", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Charles can load the file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

	})
})
