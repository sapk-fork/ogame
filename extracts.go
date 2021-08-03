package ogame

import "bytes"

// Compile time checks to ensure type satisfies Extractor interface
var _ Extractor = ExtractorV6{}
var _ Extractor = (*ExtractorV6)(nil)
var _ Extractor = ExtractorV7{}
var _ Extractor = (*ExtractorV7)(nil)

// extract universe speed from html calculation
// pageHTML := b.getPageContent(url.Values{"page": {"techtree"}, "tab": {"2"}, "techID": {"1"}})
func extractUniverseSpeed(pageHTML []byte) int64 {
	return extractUniverseSpeedV6(pageHTML)
}

func ReplaceHostname(bot *OGame, html []byte) []byte {
	serverURLBytes := []byte(bot.ServerURL())
	apiNewHostnameBytes := []byte(bot.apiNewHostname)
	escapedServerURL := bytes.Replace(serverURLBytes, []byte("/"), []byte(`\/`), -1)
	doubleEscapedServerURL := bytes.Replace(serverURLBytes, []byte("/"), []byte("\\\\\\/"), -1)
	escapedAPINewHostname := bytes.Replace(apiNewHostnameBytes, []byte("/"), []byte(`\/`), -1)
	doubleEscapedAPINewHostname := bytes.Replace(apiNewHostnameBytes, []byte("/"), []byte("\\\\\\/"), -1)
	html = bytes.Replace(html, serverURLBytes, apiNewHostnameBytes, -1)
	html = bytes.Replace(html, escapedServerURL, escapedAPINewHostname, -1)
	html = bytes.Replace(html, doubleEscapedServerURL, doubleEscapedAPINewHostname, -1)
	return html
}
