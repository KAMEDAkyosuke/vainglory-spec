package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/xmlpath.v2"
	"io"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"flag"
)

var heroName = []string{
	"Adagio",
	"Ardan",
	"Catherine",
	"Celeste",
	"Fortress",
	"Glaive",
	"Joule",
	"Koshka",
	"Krul",
	"Petal",
	"Ringo",
	"Rona",
	"SAW",
	"Skaarf",
	"Skye",
	"Taka",
	"Vox",
}

type Hero struct {
	Name         string      `json:"name"`
	HitPoints    [12]float64 `json:"hit_points"`
	HPRegen      [12]float64 `json:"hp_regen"`
	EnergyPoints [12]float64 `json:"energy_points"`
	EPRegen      [12]float64 `json:"ep_regen"`
	WeaponDamage [12]float64 `json:"weapon_damage"`
	AttackSpeed  [12]float64 `json:"attack_speed"`
	Armor        [12]float64 `json:"armor"`
	Shield       [12]float64 `json:"shield"`
	AttackRange  float64     `json:"attack_range"`
	MoveSpeed    float64     `json:"move_speed"`
}

func (hero *Hero) String() string {
	b, _ := json.MarshalIndent(hero, "", "\t")
	return string(b)
}

func NewHero(name string, body io.Reader) *Hero {
	hero := &Hero{Name: name}

	root, err := xmlpath.ParseHTML(body)
	if err != nil {
		log.Fatalf("xmlpath.ParseHTML fail : %v", err)
	}

	status := None
	path := xmlpath.MustCompile("//table/tbody/tr/td/div/span")
	iter := path.Iter(root)
	for iter.Next() {
		n := iter.Node()
		if status == None {
			switch n.String() {
			case "Hit Points (HP)":
				status = HitPoints
			case "HP Regen":
				status = HPRegen
			case "Energy Points (EP)":
				status = EnergyPoints
			case "EP Regen":
				status = EPRegen
			case "Weapon Damage":
				status = WeaponDamage
			case "Attack Speed":
				status = AttackSpeed
			case "Armor":
				status = Armor
			case "Shield":
				status = Shield
			case "Attack Range":
				status = AttackRange
			case "Move Speed":
				status = MoveSpeed
			default:
				status = None
			}
		} else {
			base, inc := parse(n.String())
			switch status {
			case HitPoints:
				for i := 0; i < 12; i++ {
					hero.HitPoints[i] = base + inc*float64(i)
					hero.HitPoints[i], _ = strconv.ParseFloat(strings.TrimSpace(fmt.Sprintf("%7.3f", hero.HitPoints[i])), 64)
				}
			case HPRegen:
				for i := 0; i < 12; i++ {
					hero.HPRegen[i] = base + inc*float64(i)
					hero.HPRegen[i], _ = strconv.ParseFloat(strings.TrimSpace(fmt.Sprintf("%7.3f", hero.HPRegen[i])), 64)
				}
			case EnergyPoints:
				for i := 0; i < 12; i++ {
					hero.EnergyPoints[i] = base + inc*float64(i)
					hero.EnergyPoints[i], _ = strconv.ParseFloat(strings.TrimSpace(fmt.Sprintf("%7.3f", hero.EnergyPoints[i])), 64)
				}
			case EPRegen:
				for i := 0; i < 12; i++ {
					hero.EPRegen[i] = base + inc*float64(i)
					hero.EPRegen[i], _ = strconv.ParseFloat(strings.TrimSpace(fmt.Sprintf("%7.3f", hero.EPRegen[i])), 64)
				}
			case WeaponDamage:
				for i := 0; i < 12; i++ {
					hero.WeaponDamage[i] = base + inc*float64(i)
					hero.WeaponDamage[i], _ = strconv.ParseFloat(strings.TrimSpace(fmt.Sprintf("%7.3f", hero.WeaponDamage[i])), 64)
				}
			case AttackSpeed:
				for i := 0; i < 12; i++ {
					hero.AttackSpeed[i] = base + inc*float64(i)
					hero.AttackSpeed[i], _ = strconv.ParseFloat(strings.TrimSpace(fmt.Sprintf("%7.3f", hero.AttackSpeed[i])), 64)
				}
			case Armor:
				for i := 0; i < 12; i++ {
					hero.Armor[i] = base + inc*float64(i)
					hero.Armor[i], _ = strconv.ParseFloat(strings.TrimSpace(fmt.Sprintf("%7.3f", hero.Armor[i])), 64)
				}
			case Shield:
				for i := 0; i < 12; i++ {
					hero.Shield[i] = base + inc*float64(i)
					hero.Shield[i], _ = strconv.ParseFloat(strings.TrimSpace(fmt.Sprintf("%7.3f", hero.Shield[i])), 64)
				}
			case AttackRange:
				hero.AttackRange = base
			case MoveSpeed:
				hero.MoveSpeed = base
			default:
				status = None
			}
			status = None
		}
	}
	return hero
}

type ParserStatus int

const (
	None ParserStatus = iota + 1
	HitPoints
	HPRegen
	EnergyPoints
	EPRegen
	WeaponDamage
	AttackSpeed
	Armor
	Shield
	AttackRange
	MoveSpeed
)

func urls(name []string) []string {
	var r []string = make([]string, len(name))
	for i, v := range name {
		r[i] = fmt.Sprintf("http://www.vaingloryfire.com/vainglory/wiki/heroes/%s/guides", strings.ToLower(v))
	}
	return r
}

func parse(s string) (float64, float64) {
	var currentStatusIsBase = true
	var base string = ""
	var inc string = ""
	for _, v := range s {
		if v == '(' {
			currentStatusIsBase = false
		} else if v == ')' {
			break
		} else if strings.Contains("0123456789.", string(v)) {
			if currentStatusIsBase {
				base += string(v)
			} else {
				inc += string(v)
			}
		}
	}

	a, err := strconv.ParseFloat(base, 64)
	if err != nil {
		// log.Fatalf("strconv.ParseFloat fail : %v", err)
		a = 0
	}
	b, err := strconv.ParseFloat(inc, 64)
	if err != nil {
		// log.Fatalf("strconv.ParseFloat fail : %v", err)
		b = 0
	}
	return a, b
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Lshortfile)

	// version
	var version = flag.String("v", "", "version")
	flag.Parse()

	type outputFormat struct {
		Version string  `json:"version"`
		Hero    []*Hero `json:"hero"`
	}

	o := &outputFormat{
		Version: *version,
		Hero:    make([]*Hero, len(heroName)),
	}
	for i, url := range urls(heroName) {
		res, err := http.Get(url)
		if err != nil {
			log.Fatalf("http.Get fail : %v", err)
		}
		defer res.Body.Close()
		o.Hero[i] = NewHero(heroName[i], res.Body)
	}

	b, _ := json.MarshalIndent(o, "", "\t")
	fmt.Printf("%v\n", string(b))
}
