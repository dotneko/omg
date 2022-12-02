# `omg`

**Onomy Manager** by [nomblocks.io](https://nomblocks.io/)

A command line tool for common user/validator interactions with the [Onomy Protocol](https://onomy.io/) blockchain.

`omg` functions as a wrapper for the `onomyd` command line tool to provide the following:

* Simple address book to store onomy/validator addresses
* Importing addresses stored in the onomy keyring
* Query balances and rewards
* Sending tokens
* Delegating and withdrawing rewards
* Automated restaking of delegator rewards

## Prerequisites

* Go v1.18+
* Locally running Onomy full node (see [Onomy Docs](https://docs.onomy.io/run-a-full-node/starting-a-full-node))
* User-owned keys already stored in the onomy keyring

## Quickstart

Quick start instructions [here](Quickstart.md)

## Installation

Clone this repo

```
git clone https://github.com/dotneko/omg.git
```

Change into the `omg` directory then run `go build .`

## Configuration

Settings can be modified in `.omgconfig.yaml`

Copy/move `.omgconfig.yaml` to home directory or `omg` binary path

## Usage

A full list of commands is shown by running `omg` with the `--help` or `-h` flag. This will also show the abbreviations for each command.

### Managing Addresses

The address book is managed using `omg address` or its alias `omg addr` and its subcommands

```
  add         Add an address to the address book
  import      Import addresses from keyring
  modify      Modify an address book entry
  rm          Delete an address book entry
  show        Show one or all addresses
```

To create an address book, you can import from the onomy keyring (default keyring is "test")

```
omg addr import
```

To specify keyring (e.g. pass), use the flag `--keyring [backend]`:
```
omg addr import --keyring pass
```

Addresses can be added using `add [alias] [address]`:
```
omg addr add user1 onomy123456789012345678901234567890123456789
```

Show list of addresses:
```
omg addr show
```

Get address for *user1*
```
omg addr show user1
```
### Validators

A list of active validators and their valoper-addresses on the chain can be queried using the `validator show` command:

```
omg validator show
```

To show commissions for a validator:

```
omg validator commissions [moniker|valoper-address]
```

### Queries

Check balances for *user1*
```
omg balances user1
```

Check rewards for *user1*
```
omg rewards user1
```

The `--all` or `-a` flag can be used to query all accounts in the address book.

```
omg balances -a
```

```
omg rewards --all
```

> N.B. Amounts displayed use the underscore `_` separator for easier reading. For amounts with a long string of digits after the decimal point (commonly when checking rewards), the digits will be truncated and shown as `._`

Raw amounts can be displayed by the `--raw` or `-r` flag.
```
omg balances -a -raw
```
### Transactions

Transactions functions are listed under `tx` command:

```
  delegate            Delegate tokens from account to validator
  restake             Withdraw rewards and restake to validator
  send                Send tokens from an account to another account/address
  withdraw-commission Withdraw commissions and rewards for validator
  withdraw-rewards    Withdraw all rewards for account
```

All transactions assume that the account `name` in the address book matches the name of the user's key in the keyring, and will fail if the onomyd cannot find the key in the keyring.

For **delegate** and **restake** commands, one can specify the `moniker` of an active validator on chain, or the validator *valoper* address. 

By default, transactions will be generated and wait for user confirmation.

The default keyring-backend is `test`, but can be modified using the flag `--keyring`. For example, to use the `pass` keyring-backend, specify `--keyring pass` when executing your command. This default could be configured in the `.omgconfig.yaml` config file.

To automate transactions, specify `--auto` or `-a`. When this flag is used, transaction prompts are automatically confirmed, therefore *be sure that the transaction is what you want to execute*. Note that this is only confirmed to work for the default keyring-backend `test`.

> N.B. **omg is a wrapper for the onomyd daemon and *DOES NOT* have access to the user's private keys/mnemonic**

#### Delegate

To delegate 100,000,000,000anom from *user1* to validator with moniker "validator1"

```
omg tx delegate user1 validator1 100000000000anom
```

Underscored '000s separators are supported for the amount

```
omg tx delegate user1 validator1 100_000_000_000anom
```

To delegate the full avaiable balance leaving a remainder, use the `--full` or `-f` flag. The default remainder is configured in the configuration file.

```
omg tx delegate user1 validator1 --full
```

To adjust the remainder amount, add the `--remainder [amount]` or `-r [amount]` flag denominated in the base denom.

```
omg tx delegate user1 validator1 --full -r 1000anom
```

> N.B. The final balance is likely to differ from the remainder set due to *auto claim rewards* being triggered by the delegation transaction.

#### Restake

The `restake` subcommand will withdraw all rewards for the account then delegate the full amount *less remainder* to the specified validator.
```
omg tx restake user1 validator1
```
Auto restake (using default remainder amount)
```
omg tx restake user1 validator1 --auto
```

Auto restake (specify remainder amount)
```
omg tx restake user1 validator1 --auto -r 1000000anom
```

#### Send

Send tokens between accounts in the address book

```
omg tx send user1 user2 1000000anom
```

Send tokens from account to external address
```
omg tx send user1 onomy1234567890123456789012345678901234567xx 1000000anom
```

#### Withdraw validator commissions and rewards

```
omg tx withdraw-commissions [validator] [moniker|valoper-address]
```
> This assumes that [validator] matches the name of the keyring and is a self-delegate of the validator

#### Withdraw rewards

Withdraw all rewards for *user1*
```
omg tx withdraw-rewards [user]
```

Automated withdraw all rewards for *user1*
```
omg tx wd user1 --auto
```

## Conversion

Convert token amount to base denom amount
```
omg convert 1nom // Returns 1000000000000000000anom
```

Convert base denom amount to token amount
```
omg convert 1000000000000000000anom   // Returns 1nom
```

```
omg convert 1_000_000_000_000_000_000anom // Returns 1nom
```
