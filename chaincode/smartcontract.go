package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Party struct {
	Municipality string `json:"Municipality"`
	Symbol       string `json:"Symbol"`
	Party        string `json:"Party"`
	Votes        int    `json:"Votes"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	parties := []Party{
		{Municipality: "Odense", Party: "Konservativ", Symbol: "C", Votes: 0},
		{Municipality: "Odense", Party: "Socialdemokratiet", Symbol: "A", Votes: 0},
		{Municipality: "Odense", Party: "Venstre", Symbol: "V", Votes: 0},
		{Municipality: "Odense", Party: "Socialistisk Folkeparti", Symbol: "F", Votes: 0},
		{Municipality: "Odense", Party: "Radikale", Symbol: "B", Votes: 0},
		{Municipality: "Odense", Party: "Enhedslisten", Symbol: "OE", Votes: 0},
		{Municipality: "Odense", Party: "Nye Borgerlige", Symbol: "D", Votes: 0},
		{Municipality: "Odense", Party: "Dansk Folkeparti", Symbol: "O", Votes: 0},
		{Municipality: "Odense", Party: "Liberaterne", Symbol: "J", Votes: 0},
		{Municipality: "Odense", Party: "Liberal Alliance", Symbol: "I", Votes: 0},
		{Municipality: "Odense", Party: "Veganerpartiet", Symbol: "G", Votes: 0},
		{Municipality: "Odense", Party: "Kristen Demokraterne", Symbol: "K", Votes: 0},
		{Municipality: "Odense", Party: "Alternativet", Symbol: "AA", Votes: 0},
	}

	for _, party := range parties {
		partyJSON, err := json.Marshal(party)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(party.Symbol, partyJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state %v", err)
		}
	}

	return nil
}

func (s *SmartContract) GetAllParties(ctx contractapi.TransactionContextInterface) ([]*Party, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	var parties []*Party
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var party Party
		err = json.Unmarshal(queryResponse.Value, &party)
		if err != nil {
			return nil, err
		}
		parties = append(parties, &party)
	}

	return parties, nil
}

func CurrentVoteCount(ctx contractapi.TransactionContextInterface) (int, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return -1, err
	}

	defer resultsIterator.Close()

	count := 0

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return -1, err
		}

		var party Party
		err = json.Unmarshal(queryResponse.Value, &party)
		if err != nil {
			continue
		}

		count += party.Votes
	}

	return count, nil
}

func (s *SmartContract) Vote(ctx contractapi.TransactionContextInterface, symbol string, votes int) (*Party, error) {
	partyJson, err := ctx.GetStub().GetState(symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	if partyJson == nil {
		return nil, fmt.Errorf("The party with a symbol of %s does not exist in Odense", symbol)
	}

	var party Party
	err = json.Unmarshal(partyJson, &party)

	if err != nil {
		return nil, err
	}

	voteCount, err := CurrentVoteCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("Unable to get current vote count, try again later please.")
	}

	if voteCount+votes > 166_955 { // https://www.kmdvalg.dk/kv/2021/K84733461.htm
		return nil, fmt.Errorf("Votes would exceed municipality eligible")
	}

	party.Votes += votes

	partyJson, jsonErr := json.Marshal(party)
	if jsonErr != nil {
		return nil, jsonErr
	}

	ctx.GetStub().PutState(symbol, partyJson)

	return &party, nil
}
