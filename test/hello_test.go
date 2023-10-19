package main
import "testing"

func TestHello(t *testing.T) {
	t.Run("say hello to a given name", func(t *testing.T) {
		got:= Hello("Chris", "English")
		want:= "Hello, Chris"
	
		assertCorrectMessage(t,got,want)
	})

	t.Run("say hello world to empty string", func(t *testing.T){
		got:= Hello("", "")
		want:= "Hello, world"

		assertCorrectMessage(t,got,want)
	})
	
	t.Run("say Hola in Spanish", func(t *testing.T) {
		got:= Hello("Elodie", "Spanish")
		want:= "Hola, Elodie"
		assertCorrectMessage(t, got,want)
	})

	t.Run("say Bonjour in French", func(t *testing.T) {
		got:= Hello("", "French")
		want:= "Bonjour, world"
		assertCorrectMessage(t, got,want)
	})
}

func assertCorrectMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}