
package akubra

import (
    "testing"
)

func TestResolveName(t *testing.T) {
    var testPairs = []struct { format string; path string}  {
        { "#", "e/info%3Afedora%2Ffedora-system%3AFedoraObject-3.0" },
        { "##", "e5/info%3Afedora%2Ffedora-system%3AFedoraObject-3.0" },
        { "#/#", "e/5/info%3Afedora%2Ffedora-system%3AFedoraObject-3.0" },
        { "#/##", "e/57/info%3Afedora%2Ffedora-system%3AFedoraObject-3.0" },
    }

    for _, w := range testPairs {
        p := Pool { Format: w.format }
        path := p.resolveName("fedora-system:FedoraObject-3.0")
        if path != w.path {
            t.Error(path)
        }
    }
}

func TestGuessFormat(t *testing.T) {
    // These test paths are wildly wrong. Fix them!
    var testPairs = []struct { path, format string } {
        {"/Users/dbrower/Documents/work/hydra-tutorial/hydra-tut/jetty/fedora/default/data/objectStore", "##"},
        {"/Users/dbrower/Documents/work/hydra-tutorial/hydra-tut/jetty/fedora/default/data", ""},
        {"/Users/dbrower/Documents/work/go/src/github.com/dbrower/fedoro/akubra/test-dir", "#/##"},
    }

    for _, w := range testPairs {
        guess := GuessFormat(w.path)
        if guess != w.format {
            t.Errorf("For %s expected %s, got %s", w.path, w.format, guess)
        }
    }
}
