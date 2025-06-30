package config

import (
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"

	"github.com/oschwald/geoip2-golang"
)

// List of denied country ISO codes
var DeniedCountries = map[string]bool{
	"RU": true, // Russia
	"KP": true, // North Korea
	"IR": true, // Iran
	"SY": true, // Syria
	"BY": true, // Belarus
}

// Get the public IP of the server
func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
}

// Check if the server is running in a denied country
func checkServerCountryBlock(db *geoip2.Reader) {
	ipStr, err := getPublicIP()
	if err != nil {
		log.Fatalf("Could not determine public IP: %v", err)
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		log.Fatalf("Invalid public IP: %s", ipStr)
	}
	record, err := db.Country(ip)
	if err != nil {
		log.Fatalf("Could not determine server country: %v", err)
	}
	if DeniedCountries[record.Country.IsoCode] {
		log.Fatalf("Server is running in a denied country (%s). Exiting.", record.Country.IsoCode)
	}
}

// Try to open the GeoIP database from embedded FS, then from disk as fallback
func OpenGeoIPDB(embeddedFS fs.FS, localPath string) (*geoip2.Reader, error) {
	if embeddedFS != nil {
		f, err := embeddedFS.Open("data/geo/GeoLite2-Country.mmdb")
		if err == nil {
			defer f.Close()
			data, err := io.ReadAll(f)
			if err == nil {
				return geoip2.FromBytes(data)
			}
		}
	}
	return geoip2.Open(localPath)
}

func (c *Config) CheckServerCountryBlock() {
	db, err := OpenGeoIPDB(c.Database.GeoIPFS, c.Database.Path+"/geo/GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatalf("Failed to open GeoIP database: %v", err)
	}
	defer db.Close()

	checkServerCountryBlock(db)
}

func (c *Config) GetGeoIPDB() *geoip2.Reader {
	db, err := OpenGeoIPDB(c.Database.GeoIPFS, c.Database.Path+"/geo/GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatalf("Failed to open GeoIP database: %v", err)
	}
	return db
}
