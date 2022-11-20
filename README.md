# `omg`

**Onomy Manager** by [nomblocks.io](https://nomblocks.io/)

A command line tool for common user/validator interactions with the [Onomy Protocol[(https://onomy.io/)] blockchain.

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

## Installation

Clone this repo

```
git clone https://github.com/dotneko/omg.git
```

Change into the `omg` directory then run `go build .`

## Configuration

Settings can be modified in `.omgconfig.yaml`

Ensure `.omgconfig.yaml` in home directory or binary path

## Usage

A full list of commands is shown by running `omg` without any flags

### Managing Addresses

The address book is managed using `omg addr` and its subcommands

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

Validator addresses can also be added using `add [alias] [valoper address]`:

```
omg addr add validator1 onomyvaloper123456789012345678901234567890123456789
```

Show list of addresses:
```
omg addr show
```

Get address for *user1*
```
omg addr show user1
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
omg rewards --all
```

### Transactions

Transactions functions are listed under `tx` command:

```
  delegate    Delegate tokens from account to validator
  restake     Restake rewards for account to validator
  send        Send tokens from an account to another account/address
  wdrewards   Withdraw all rewards
```

All transactions assume that the account `name` in the address book matches the name of the user's key in the keyring, and will fail if the onomyd cannot find the key in the keyring.

By default, transactions will be generated and wait for user confirmation.

The default keyring-backend is `test`, but can be modified using the flag `--keyring`. For example, to use the `pass` keyring-backend, specify `--keyring pass` when executing your command.

To automate transactions, specify `--auto` or `-a`. When this flag is used, transaction prompts are automatically confirmed, therefore be sure that the transaction is what you want to execute. Note that this is only confirmed to work for the default keyring-backend `test`.

> N.B. **omg is a wrapper for the onomyd daemon and *DOES NOT* have access to the user's private keys/mnemonic**

#### Delegate

To delegate 100,000,000,000anom from *user1* to *validator1*
(If no amount specified, user will be prompted for the amount)

```
omg tx delegate user1 validator1 100000000000anom
```

To delegate the full avaiable balance leaving a remainder, use the `--full` or `-f` flag

```
omg tx delegate user1 validator1 --full
```

To adjust the remainder amount, user the `--remainder [amount]` or `-r [amount]` flag
```
omg tx delegate user1 validator1 --full -r 1000anom
```

#### Restake

The `restake` subcommand will withdraw all rewards for the account then delegate the full amount *less remainder*
to the specified validator.
```
omg tx restake user1 validator1
```
Auto restake (bypassing confirmations)
```
omg tx restake user1 validator1 --auto -r 1000000anom
```

#### Send

Send tokens between accounts in the address book
(If no amount specified, user will be prompted for the amount)

```
omg tx send user1 user2 1000000anom
```

#### Withdraw rewards

Automated withdraw all rewards for *user1*
```
omg tx wdrewards user1 --auto
```

## Conversion (Imprecise)

> **The current implementation only provides a rough conversion and is IMPRECISE**

Convert token amount to denom amount
```
omg convert 1  // Returns 1000000000000000000anom
```
Convert denom amount to token amount
```
omg convert 1000000000000000000anom // Returns 1 nom
```