package ogame

import (
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

// Priorities
const (
	Low       = 1
	Normal    = 2
	Important = 3
	Critical  = 4
)

// Prioritize ...
type Prioritize struct {
	bot          *OGame
	initiator    string
	name         string
	taskIsDoneCh chan struct{}
	isTx         int32
}

// SetInitiator ...
func (b *Prioritize) SetInitiator(initiator string) Prioritizable {
	b.initiator = initiator
	return b
}

// Begin a new transaction. "Done" must be called to release the lock.
func (b *Prioritize) Begin() Prioritizable {
	return b.BeginNamed("Tx")
}

// BeginNamed begins a new transaction with a name. "Done" must be called to release the lock.
func (b *Prioritize) BeginNamed(name string) Prioritizable {
	if name == "" {
		name = "Tx"
	}
	return b.begin(name)
}

// Done terminate the transaction, release the lock.
func (b *Prioritize) Done() {
	b.done()
}

func (b *Prioritize) begin(name string) *Prioritize {
	if atomic.AddInt32(&b.isTx, 1) == 1 {
		if b.initiator != "" {
			b.name = b.initiator + ":"
		}
		b.name += name
		b.bot.botLock(b.name)
	}
	return b
}

func (b *Prioritize) done() {
	if atomic.AddInt32(&b.isTx, -1) == 0 {
		defer close(b.taskIsDoneCh)
		b.bot.botUnlock(b.name)
	}
}

// Tx locks the bot during the transaction and ensure the lock is released afterward
func (b *Prioritize) Tx(clb func(Prioritizable) error) error {
	tx := b.Begin()
	defer tx.Done()
	err := clb(tx)
	return err
}

// LoginWithBearerToken to ogame server reusing existing token
// Returns either or not the bot logged in using the existing cookies
func (b *Prioritize) LoginWithBearerToken(token string) (bool, error) {
	b.begin("LoginWithBearerToken")
	defer b.done()
	return b.bot.wrapLoginWithBearerToken(token)
}

// LoginWithExistingCookies to ogame server reusing existing cookies
// Returns either or not the bot logged in using the existing cookies
func (b *Prioritize) LoginWithExistingCookies() (bool, error) {
	b.begin("LoginWithExistingCookies")
	defer b.done()
	return b.bot.wrapLoginWithExistingCookies()
}

// Login to ogame server
// Can fails with BadCredentialsError
func (b *Prioritize) Login() error {
	b.begin("Login")
	defer b.done()
	return b.bot.wrapLogin()
}

// Logout the bot from ogame server
func (b *Prioritize) Logout() {
	b.begin("Logout")
	defer b.done()
	b.bot.logout()
}

// GetAlliancePageContent gets the html for a specific ogame page
func (b *Prioritize) GetAlliancePageContent(vals url.Values) ([]byte, error) {
	b.begin("GetAlliancePageContent")
	defer b.done()
	return b.bot.getAlliancePageContent(vals)
}

// GetPageContent gets the html for a specific ogame page
func (b *Prioritize) GetPageContent(vals url.Values) ([]byte, error) {
	b.begin("GetPageContent")
	defer b.done()
	return b.bot.getPageContent(vals)
}

// PostPageContent make a post request to ogame server
// This is useful when simulating a web browser
func (b *Prioritize) PostPageContent(vals, payload url.Values) ([]byte, error) {
	b.begin("PostPageContent")
	defer b.done()
	return b.bot.postPageContent(vals, payload)
}

// IsUnderAttack returns true if the user is under attack, false otherwise
func (b *Prioritize) IsUnderAttack() (bool, error) {
	b.begin("IsUnderAttack")
	defer b.done()
	return b.bot.isUnderAttack()
}

// GetPlanets returns the user planets
func (b *Prioritize) GetPlanets() []Planet {
	b.begin("GetPlanets")
	defer b.done()
	return b.bot.getPlanets()
}

// GetPlanet gets infos for planetID
// Fails if planetID is invalid
func (b *Prioritize) GetPlanet(v interface{}) (Planet, error) {
	b.begin("GetPlanet")
	defer b.done()
	return b.bot.getPlanet(v)
}

// GetMoons returns the user moons
func (b *Prioritize) GetMoons() []Moon {
	b.begin("GetMoons")
	defer b.done()
	return b.bot.getMoons()
}

// GetMoon gets infos for moonID
func (b *Prioritize) GetMoon(v interface{}) (Moon, error) {
	b.begin("GetMoon")
	defer b.done()
	return b.bot.getMoon(v)
}

// GetCelestials get the player's planets & moons
func (b *Prioritize) GetCelestials() ([]Celestial, error) {
	b.begin("GetCelestials")
	defer b.done()
	return b.bot.getCelestials()
}

// RecruitOfficer recruit an officer.
// Typ 2: Commander, 3: Admiral, 4: Engineer, 5: Geologist, 6: Technocrat
// Days: 7 or 90
func (b *Prioritize) RecruitOfficer(typ, days int64) error {
	b.begin("RecruitOfficer")
	defer b.done()
	return b.bot.recruitOfficer(typ, days)
}

// Abandon a planet. Warning: this is irreversible
func (b *Prioritize) Abandon(v interface{}) error {
	b.begin("Abandon")
	defer b.done()
	return b.bot.abandon(v)
}

// GetCelestial get the player's planet/moon using the coordinate
func (b *Prioritize) GetCelestial(v interface{}) (Celestial, error) {
	b.begin("GetCelestial")
	defer b.done()
	return b.bot.getCelestial(v)
}

// ServerTime returns server time
// Timezone is OGT (OGame Time zone)
func (b *Prioritize) ServerTime() time.Time {
	b.begin("ServerTime")
	defer b.done()
	return b.bot.serverTime()
}

// GetUserInfos gets the user information
func (b *Prioritize) GetUserInfos() UserInfos {
	b.begin("GetUserInfos")
	defer b.done()
	return b.bot.getUserInfos()
}

// SendMessage sends a message to playerID
func (b *Prioritize) SendMessage(playerID int64, message string) error {
	b.begin("SendMessage")
	defer b.done()
	return b.bot.sendMessage(playerID, message, true)
}

// SendMessageAlliance sends a message to associationID
func (b *Prioritize) SendMessageAlliance(associationID int64, message string) error {
	b.begin("SendMessageAlliance")
	defer b.done()
	return b.bot.sendMessage(associationID, message, false)
}

// GetFleets get the player's own fleets activities
func (b *Prioritize) GetFleets(opts ...Option) ([]Fleet, Slots) {
	b.begin("GetFleets")
	defer b.done()
	return b.bot.getFleets(opts...)
}

// GetFleetsFromEventList get the player's own fleets activities
func (b *Prioritize) GetFleetsFromEventList() []Fleet {
	b.begin("GetFleets")
	defer b.done()
	return b.bot.getFleetsFromEventList()
}

// CancelFleet cancel a fleet
func (b *Prioritize) CancelFleet(fleetID FleetID) error {
	b.begin("CancelFleet")
	defer b.done()
	return b.bot.cancelFleet(fleetID)
}

// GetAttacks get enemy fleets attacking you
func (b *Prioritize) GetAttacks(opts ...Option) ([]AttackEvent, error) {
	b.begin("GetAttacks")
	defer b.done()
	return b.bot.getAttacks(opts...)
}

// GalaxyInfos get information of all planets and moons of a solar system
func (b *Prioritize) GalaxyInfos(galaxy, system int64, options ...Option) (SystemInfos, error) {
	b.begin("GalaxyInfos")
	defer b.done()
	return b.bot.galaxyInfos(galaxy, system, options...)
}

// GetResourceSettings gets the resources settings for specified planetID
func (b *Prioritize) GetResourceSettings(planetID PlanetID, options ...Option) (ResourceSettings, error) {
	b.begin("GetResourceSettings")
	defer b.done()
	return b.bot.getResourceSettings(planetID, options...)
}

// SetResourceSettings set the resources settings on a planet
func (b *Prioritize) SetResourceSettings(planetID PlanetID, settings ResourceSettings) error {
	b.begin("SetResourceSettings")
	defer b.done()
	return b.bot.setResourceSettings(planetID, settings)
}

// GetResourcesBuildings gets the resources buildings levels
func (b *Prioritize) GetResourcesBuildings(celestialID CelestialID, options ...Option) (ResourcesBuildings, error) {
	b.begin("GetResourcesBuildings")
	defer b.done()
	return b.bot.getResourcesBuildings(celestialID, options...)
}

// GetDefense gets all the defenses units information of a planet
// Fails if planetID is invalid
func (b *Prioritize) GetDefense(celestialID CelestialID, options ...Option) (DefensesInfos, error) {
	b.begin("GetDefense")
	defer b.done()
	return b.bot.getDefense(celestialID, options...)
}

// GetShips gets all ships units information of a planet
func (b *Prioritize) GetShips(celestialID CelestialID, options ...Option) (ShipsInfos, error) {
	b.begin("GetShips")
	defer b.done()
	return b.bot.getShips(celestialID, options...)
}

// GetFacilities gets all facilities information of a planet
func (b *Prioritize) GetFacilities(celestialID CelestialID, options ...Option) (Facilities, error) {
	b.begin("GetFacilities")
	defer b.done()
	return b.bot.getFacilities(celestialID, options...)
}

// GetProduction get what is in the production queue.
// (ships & defense being built)
func (b *Prioritize) GetProduction(celestialID CelestialID) ([]Quantifiable, int64, error) {
	b.begin("GetProduction")
	defer b.done()
	return b.bot.getProduction(celestialID)
}

// GetCachedResearch gets the player cached researches information
func (b *Prioritize) GetCachedResearch() Researches {
	b.begin("GetCachedResearch")
	defer b.done()
	return b.bot.getCachedResearch()
}

// GetResearch gets the player researches information
func (b *Prioritize) GetResearch() Researches {
	b.begin("GetResearch")
	defer b.done()
	return b.bot.getResearch()
}

// GetSlots gets the player current and total slots information
func (b *Prioritize) GetSlots() Slots {
	b.begin("GetSlots")
	defer b.done()
	return b.bot.getSlots()
}

// Build builds any ogame objects (building, technology, ship, defence)
func (b *Prioritize) Build(celestialID CelestialID, id ID, nbr int64) error {
	b.begin("Build")
	defer b.done()
	return b.bot.build(celestialID, id, nbr)
}

// TearDown tears down any ogame building
func (b *Prioritize) TearDown(celestialID CelestialID, id ID) error {
	b.begin("TearDown")
	defer b.done()
	return b.bot.tearDown(celestialID, id)
}

// BuildCancelable builds any cancelable ogame objects (building, technology)
func (b *Prioritize) BuildCancelable(celestialID CelestialID, id ID) error {
	b.begin("BuildCancelable")
	defer b.done()
	return b.bot.buildCancelable(celestialID, id)
}

// BuildProduction builds any line production ogame objects (ship, defence)
func (b *Prioritize) BuildProduction(celestialID CelestialID, id ID, nbr int64) error {
	b.begin("BuildProduction")
	defer b.done()
	return b.bot.buildProduction(celestialID, id, nbr)
}

// BuildBuilding ensure what is being built is a building
func (b *Prioritize) BuildBuilding(celestialID CelestialID, buildingID ID) error {
	b.begin("BuildBuilding")
	defer b.done()
	return b.bot.buildBuilding(celestialID, buildingID)
}

// BuildDefense builds a defense unit
func (b *Prioritize) BuildDefense(celestialID CelestialID, defenseID ID, nbr int64) error {
	b.begin("BuildDefense")
	defer b.done()
	return b.bot.buildDefense(celestialID, defenseID, nbr)
}

// BuildShips builds a ship unit
func (b *Prioritize) BuildShips(celestialID CelestialID, shipID ID, nbr int64) error {
	b.begin("BuildShips")
	defer b.done()
	return b.bot.buildShips(celestialID, shipID, nbr)
}

// ConstructionsBeingBuilt returns the building & research being built, and the time remaining (secs)
func (b *Prioritize) ConstructionsBeingBuilt(celestialID CelestialID) (ID, int64, ID, int64) {
	b.begin("ConstructionsBeingBuilt")
	defer b.done()
	return b.bot.constructionsBeingBuilt(celestialID)
}

// CancelBuilding cancel the construction of a building on a specified planet
func (b *Prioritize) CancelBuilding(celestialID CelestialID) error {
	b.begin("CancelBuilding")
	defer b.done()
	return b.bot.cancelBuilding(celestialID)
}

// CancelResearch cancel the research
func (b *Prioritize) CancelResearch(celestialID CelestialID) error {
	b.begin("CancelResearch")
	defer b.done()
	return b.bot.cancelResearch(celestialID)
}

// BuildTechnology ensure that we're trying to build a technology
func (b *Prioritize) BuildTechnology(celestialID CelestialID, technologyID ID) error {
	b.begin("BuildTechnology")
	defer b.done()
	return b.bot.buildTechnology(celestialID, technologyID)
}

// GetResources gets user resources
func (b *Prioritize) GetResources(celestialID CelestialID) (Resources, error) {
	b.begin("GetResources")
	defer b.done()
	return b.bot.getResources(celestialID)
}

// GetResourcesDetails gets user resources
func (b *Prioritize) GetResourcesDetails(celestialID CelestialID) (ResourcesDetails, error) {
	b.begin("GetResourcesDetails")
	defer b.done()
	return b.bot.getResourcesDetails(celestialID)
}

// GetTechs gets a celestial supplies/facilities/ships/researches
func (b *Prioritize) GetTechs(celestialID CelestialID) (ResourcesBuildings, Facilities, ShipsInfos, DefensesInfos, Researches, error) {
	b.begin("GetTechs")
	defer b.done()
	return b.bot.getTechs(celestialID)
}

// SendFleet sends a fleet
func (b *Prioritize) SendFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate,
	mission MissionID, resources Resources, holdingTime, unionID int64) (Fleet, error) {
	b.begin("SendFleet")
	defer b.done()
	return b.bot.sendFleet(celestialID, ships, speed, where, mission, resources, holdingTime, unionID, false)
}

// EnsureFleet either sends all the requested ships or fail
func (b *Prioritize) EnsureFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate,
	mission MissionID, resources Resources, holdingTime, unionID int64) (Fleet, error) {
	b.begin("EnsureFleet")
	defer b.done()
	return b.bot.sendFleet(celestialID, ships, speed, where, mission, resources, holdingTime, unionID, true)
}

// DestroyRockets destroys anti-ballistic & inter-planetary missiles
func (b *Prioritize) DestroyRockets(planetID PlanetID, abm, ipm int64) error {
	b.begin("DestroyRockets")
	defer b.done()
	return b.bot.destroyRockets(planetID, abm, ipm)
}

// SendIPM sends IPM
func (b *Prioritize) SendIPM(planetID PlanetID, coord Coordinate, nbr int64, priority ID) (int64, error) {
	b.begin("SendIPM")
	defer b.done()
	return b.bot.sendIPM(planetID, coord, nbr, priority)
}

// GetCombatReportSummaryFor gets the latest combat report for a given coordinate
func (b *Prioritize) GetCombatReportSummaryFor(coord Coordinate) (CombatReportSummary, error) {
	b.begin("GetCombatReportSummaryFor")
	defer b.done()
	return b.bot.getCombatReportFor(coord)
}

// GetEspionageReportFor gets the latest espionage report for a given coordinate
func (b *Prioritize) GetEspionageReportFor(coord Coordinate) (EspionageReport, error) {
	b.begin("GetEspionageReportFor")
	defer b.done()
	return b.bot.getEspionageReportFor(coord)
}

// GetEspionageReportMessages gets the summary of each espionage reports
func (b *Prioritize) GetEspionageReportMessages() ([]EspionageReportSummary, error) {
	b.begin("GetEspionageReportMessages")
	defer b.done()
	return b.bot.getEspionageReportMessages()
}

// CollectAllMarketplaceMessages collect all marketplace messages
func (b *Prioritize) CollectAllMarketplaceMessages() error {
	b.begin("CollectAllMarketplaceMessages")
	defer b.done()
	return b.bot.collectAllMarketplaceMessages()
}

// CollectMarketplaceMessage collect marketplace message
func (b *Prioritize) CollectMarketplaceMessage(msg MarketplaceMessage) error {
	b.begin("CollectMarketplaceMessage")
	defer b.done()
	_, err := b.bot.collectMarketplaceMessage(msg, "")
	return err
}

// GetExpeditionMessages gets the expedition messages
func (b *Prioritize) GetExpeditionMessages() ([]ExpeditionMessage, error) {
	b.begin("GetExpeditionMessages")
	defer b.done()
	return b.bot.getExpeditionMessages()
}

// GetExpeditionMessageAt gets the expedition message for time t
func (b *Prioritize) GetExpeditionMessageAt(t time.Time) (ExpeditionMessage, error) {
	b.begin("GetExpeditionMessageAt")
	defer b.done()
	return b.bot.getExpeditionMessageAt(t)
}

// GetEspionageReport gets a detailed espionage report
func (b *Prioritize) GetEspionageReport(msgID int64) (EspionageReport, error) {
	b.begin("GetEspionageReport")
	defer b.done()
	return b.bot.getEspionageReport(msgID)
}

// DeleteMessage deletes a message from the mail box
func (b *Prioritize) DeleteMessage(msgID int64) error {
	b.begin("DeleteMessage")
	defer b.done()
	return b.bot.deleteMessage(msgID)
}

// DeleteAllMessagesFromTab ...
func (b *Prioritize) DeleteAllMessagesFromTab(tabID int64) error {
	b.begin("DeleteAllMessagesFromTab")
	defer b.done()
	return b.bot.deleteAllMessagesFromTab(tabID)
}

// GetResourcesProductions gets the planet resources production
func (b *Prioritize) GetResourcesProductions(planetID PlanetID) (Resources, error) {
	b.begin("GetResourcesProductions")
	defer b.done()
	return b.bot.getResourcesProductions(planetID)
}

// GetResourcesProductionsLight gets the planet resources production
func (b *Prioritize) GetResourcesProductionsLight(resBuildings ResourcesBuildings, researches Researches,
	resSettings ResourceSettings, temp Temperature) Resources {
	b.begin("GetResourcesProductionsLight")
	defer b.done()
	return getResourcesProductionsLight(resBuildings, researches, resSettings, temp, b.bot.serverData.Speed)
}

// FlightTime calculate flight time and fuel needed
func (b *Prioritize) FlightTime(origin, destination Coordinate, speed Speed, ships ShipsInfos, missionID MissionID) (secs, fuel int64) {
	b.begin("FlightTime")
	defer b.done()
	researches := b.bot.getCachedResearch()
	return CalcFlightTime(origin, destination, b.bot.serverData.Galaxies, b.bot.serverData.Systems,
		b.bot.serverData.DonutGalaxy, b.bot.serverData.DonutSystem, b.bot.serverData.GlobalDeuteriumSaveFactor,
		float64(speed)/10, GetFleetSpeedForMission(b.bot.IsV81(), b.bot.serverData, missionID), ships, researches, b.bot.characterClass)
}

// Phalanx scan a coordinate from a moon to get fleets information
// IMPORTANT: My account was instantly banned when I scanned an invalid coordinate.
// IMPORTANT: This function DOES validate that the coordinate is a valid planet in range of phalanx
// 			  and that you have enough deuterium.
func (b *Prioritize) Phalanx(moonID MoonID, coord Coordinate) ([]Fleet, error) {
	b.begin("Phalanx")
	defer b.done()
	return b.bot.getPhalanx(moonID, coord)
}

// UnsafePhalanx same as Phalanx but does not perform any input validation.
func (b *Prioritize) UnsafePhalanx(moonID MoonID, coord Coordinate) ([]Fleet, error) {
	b.begin("Phalanx")
	defer b.done()
	return b.bot.getUnsafePhalanx(moonID, coord)
}

// JumpGate sends ships through a jump gate.
func (b *Prioritize) JumpGate(origin, dest MoonID, ships ShipsInfos) (bool, int64, error) {
	b.begin("JumpGate")
	defer b.done()
	return b.bot.executeJumpGate(origin, dest, ships)
}

// JumpGateDestinations returns available destinations for jump gate.
func (b *Prioritize) JumpGateDestinations(origin MoonID) ([]MoonID, int64, error) {
	b.begin("JumpGateDestinations")
	defer b.done()
	return b.bot.jumpGateDestinations(origin)
}

// BuyOfferOfTheDay buys the offer of the day.
func (b *Prioritize) BuyOfferOfTheDay() error {
	b.begin("BuyOfferOfTheDay")
	defer b.done()
	return b.bot.buyOfferOfTheDay()
}

// CreateUnion creates a union
func (b *Prioritize) CreateUnion(fleet Fleet, users []string) (int64, error) {
	b.begin("CreateUnion")
	defer b.done()
	return b.bot.createUnion(fleet, users)
}

// HeadersForPage gets the headers for a specific ogame page
func (b *Prioritize) HeadersForPage(url string) (http.Header, error) {
	b.begin("HeadersForPage")
	defer b.done()
	return b.bot.headersForPage(url)
}

// GetEmpire (Commander only)
func (b *Prioritize) GetEmpire(celestialType CelestialType) ([]EmpireCelestial, error) {
	b.begin("GetEmpire")
	defer b.done()
	return b.bot.getEmpire(celestialType)
}

// GetEmpireJSON retrieves JSON from Empire page (Commander only).
func (b *Prioritize) GetEmpireJSON(nbr int64) (interface{}, error) {
	b.begin("GetEmpireJSON")
	defer b.done()
	return b.bot.getEmpireJSON(nbr)
}

// GetAuction ...
func (b *Prioritize) GetAuction() (Auction, error) {
	b.begin("GetAuction")
	defer b.done()
	return b.bot.getAuction(CelestialID(0))
}

// DoAuction ...
func (b *Prioritize) DoAuction(bid map[CelestialID]Resources) error {
	b.begin("DoAuction")
	defer b.done()
	return b.bot.doAuction(CelestialID(0), bid)
}

// Highscore ...
func (b *Prioritize) Highscore(category, typ, page int64) (Highscore, error) {
	b.begin("Highscore")
	defer b.done()
	return b.bot.highscore(category, typ, page)
}

// GetAllResources ...
func (b *Prioritize) GetAllResources() (map[CelestialID]Resources, error) {
	b.begin("GetAllResources")
	defer b.done()
	return b.bot.getAllResources()
}

// GetDMCosts returns fast build with DM information
func (b *Prioritize) GetDMCosts(celestialID CelestialID) (DMCosts, error) {
	b.begin("GetDMCosts")
	defer b.done()
	return b.bot.getDMCosts(celestialID)
}

// UseDM use dark matter to fast build
func (b *Prioritize) UseDM(typ string, celestialID CelestialID) error {
	b.begin("UseDM")
	defer b.done()
	return b.bot.useDM(typ, celestialID)
}

// GetItems get all items information
func (b *Prioritize) GetItems(celestialID CelestialID) ([]Item, error) {
	b.begin("GetItems")
	defer b.done()
	return b.bot.getItems(celestialID)
}

// GetActiveItems ...
func (b *Prioritize) GetActiveItems(celestialID CelestialID) ([]ActiveItem, error) {
	b.begin("GetActiveItems")
	defer b.done()
	return b.bot.getActiveItems(celestialID)
}

// ActivateItem activate an item
func (b *Prioritize) ActivateItem(ref string, celestialID CelestialID) error {
	b.begin("ActivateItem")
	defer b.done()
	return b.bot.activateItem(ref, celestialID)
}

// BuyMarketplace buy an item on the marketplace
func (b *Prioritize) BuyMarketplace(itemID int64, celestialID CelestialID) error {
	b.begin("BuyMarketplace")
	defer b.done()
	return b.bot.buyMarketplace(itemID, celestialID)
}

// OfferSellMarketplace ...
func (b *Prioritize) OfferSellMarketplace(itemID interface{}, quantity, priceType, price, priceRange int64, celestialID CelestialID) error {
	b.begin("OfferSellMarketplace")
	defer b.done()
	return b.bot.offerMarketplace(4, itemID, quantity, priceType, price, priceRange, celestialID)
}

// OfferBuyMarketplace ...
func (b *Prioritize) OfferBuyMarketplace(itemID interface{}, quantity, priceType, price, priceRange int64, celestialID CelestialID) error {
	b.begin("OfferBuyMarketplace")
	defer b.done()
	return b.bot.offerMarketplace(3, itemID, quantity, priceType, price, priceRange, celestialID)
}
