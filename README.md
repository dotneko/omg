# `omg`

**Onomy Manager** by nomblocks.io

A command line tool for common user / validator interactions with the Onomy Protocol blockchain.

`omg` functions as a wrapper for the `onomyd` command line tool to provide the following:

* Simple address book to store onomy/validator addresses
* Importing addresses stored in the onomy keyring
* Conversion between the token `nom` <-> `anom` denom amounts
* Query balances and rewards
* Single command auto-restaking delegator rewards

## Prerequisites

* Go v1.18
* Locally running Onomy full node

## Installation

Clone this repo

```
git clone https://github.com/dotneko/omg.git
```

Change into the `omg/cmd/omg` directory then run `go build .`
```
cd omg/cmd/omg
go build .
```

## Configuration

Set location for wallet data using environmental variable:

```
export OMG_FILENAME="/home/user/.omg.json"
```

## Usage

A full list of commands is shown by running `omg` without any flags

### Managing Addresses

To create address book, you can import from the onomy keyring (default keyring is "test")

```
omg -import
```

To specify keyring (e.g. pass), use the flag `-keyring [backend]`:
```
omg -import -keyring pass
```

Addresses can be added using `-add [alias] [address]`:
```
omg -add some_alias onomy12345678901234567890123456789
```

Show list of addresses and their aliases:
```
omg -list
```
### Queries

Check balance for *alias*
```
omg -balances alias
```

Check rewards for *alias*
```
omg -rewards alias
```

### Transactions

#### Send
```
omg -send from_alias to_alias
```

#### Withdraw all rewards
```
omg -wdall alias
```

#### Delegate

Delegate amount to be input at prompt after executing
```
omg -delegate alias validator_alias
```

To specify delegate amount
```
omg -delegate alias validator_alias 1000000000000anom
```

To specify remainder amount, use a negative value
```
omg -delegate alias validator_alias -100000000anom
```

Values without the denom suffix is treated as the token amount, e.g.
```
omg -delegate alias validator_alias -1
```
would delegate available balance leaving a remainder of 1 nom (1000000000000000000anom)

#### Restake

The following will withdraw all rewards for the account, then delegate *entire balance less 1 nom*
```
omg -restake alias validator_alias
```

## Conversion

Convert token amount to denom amount
```
omg -convt 1  // Returns 1000000000000000000anom
```
Convert denom amount to token amount
```
omg -convd 1000000000000000000anom // Returns 1 nom
```