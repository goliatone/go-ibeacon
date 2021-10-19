// package beacon
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/go-ble/ble"
	// "github.com/go-ble/ble/darwin"
	"github.com/go-ble/ble/examples/lib"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/pkg/errors"
)

//VERSION Generated via ld flags
var VERSION string

//BUILD_DATE Generated via ld flags
var BUILD_DATE string

//IBeacon - Structure to hold iBeacon data
type IBeacon struct {
	serviceList []*ble.Service
	Device      ble.Device

	UUID       string
	Name       string
	Major      uint16
	Minor      uint16
	PowerLevel int8
}

//NewIBeacon - Create new IBeacon instance
func NewIBeacon(uuid string, name string, powerLevel int8) *IBeacon {
	device, err := dev.NewDevice("default")
	// device, err := darwin.NewDevice()
	if err != nil {
		log.Fatalf("Unable to unitialize new device: %v", err)
	}

	return &IBeacon{
		Device:     device,
		UUID:       uuid,
		Name:       name,
		Major:      1,
		Minor:      1,
		PowerLevel: powerLevel,
	}
}

//SetiBeaconVersion - Update beacon version
func (beacon *IBeacon) SetiBeaconVersion(major, minor uint16) {
	beacon.Major = major
	beacon.Minor = minor
}

//AddBatteryService - advertise battery
func (beacon *IBeacon) AddBatteryService() {
	service := ble.NewService(ble.BatteryUUID)
	beacon.serviceList = append(beacon.serviceList, service)
}

//AddCountService - counter
func (beacon *IBeacon) AddCountService() {
	service := ble.NewService(lib.TestSvcUUID)
	service.AddCharacteristic(lib.NewCountChar())
	service.AddCharacteristic(lib.NewEchoChar())
}

//Advertise - Start advertisement
func (beacon *IBeacon) Advertise(duration uint64) error {
	ble.SetDefaultDevice(beacon.Device)

	for _, service := range beacon.serviceList {
		if err := ble.AddService(service); err != nil {
			log.Fatalf("Unable to add service: %v", err)
		}
	}

	if duration == 0 {
		fmt.Println("Advertising until quitting...")
	} else {
		fmt.Printf("Advertising for %s...\n", time.Duration(duration)*time.Second)
	}

	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), time.Duration(duration)))
	return ble.AdvertiseNameAndServices(ctx, beacon.Name)
}

func main() {

	// EC9E84F8-87D8-498B-8B0C-9EF8D3AA94C7
	// 01020304050607080910111213141516
	uuid := flag.String("uuid", "EC9E84F8-87D8-498B-8B0C-9EF8D3AA94C7", "UUID used to advertise the beacon")
	major := flag.Int("major", 1, "iBeacon major version")
	minor := flag.Int("minor", 1, "iBeacon minor version")
	name := flag.String("name", "rpibeacon", "iBeacon name")
	power := flag.Int("power-level", -60, "iBeacon power level")
	duration := flag.Int("duration", 0, "Advertise for duration, 0 is for ever")

	flag.Parse()

	ibeacon := NewIBeacon(*uuid, *name, int8(*power))
	ibeacon.SetiBeaconVersion(uint16(*major), uint16(*minor))
	// ibeacon.AddCountService()
	// ibeacon.AddBatteryService()
	err := ibeacon.Advertise(uint64(*duration))

	//TODO: Handle signals, ensure we have a way to quit!
	checkError(err)
}

func checkError(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Println("done")
	case context.Canceled:
		fmt.Println("canceled")
	default:
		log.Fatalf(err.Error())
	}
}
