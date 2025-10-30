// pkg/twiml/generator.go
package twiml

import "fmt"

func SayGather(text, action string) string {
	// Removed: <?xml version="1.0" encoding="UTF-8"?>
	return fmt.Sprintf(`<Response>
  <Say voice="man">%s</Say>
  <Gather input="speech" action="%s" method="POST" speechTimeout="auto">
    <Say>Please speak now.</Say>
  </Gather>
  <Say>I didn't hear anything. Goodbye.</Say>
</Response>`, text, action)
}

func SayAndHangup(text string) string {
	// Removed: <?xml version="1.0" encoding="UTF-8"?>
	return fmt.Sprintf(`<Response>
  <Say voice="man">%s</Say>
  <Hangup/>
</Response>`, text)
}

