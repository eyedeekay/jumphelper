# jumphelper

This is a local service which you can use to query i2p jump services but only as
needed, with rate-limiting and a place to experiment with privacy-enhancing
measures, like limiting the overall and remotely deducible uniqueness of a local
addressbook while not compromising it's ability to be defined by the local user.

If you want to test it, that's fine, but it's highly experimental.

## GOALS:

### Service/Library goals

  * Make it more difficult for an eepSite with embedded resources from remote
  sites to reliably determine the long-term contents of a user's address book
  by rate-limiting requests to the local addressbook.(Works, but I don't know
  what to do about the defaults.)
  * Make it easier for third-party, non-core i2p applications to create their
  own, per-application address books. This will allow them to construct address
  books based on what they need to contain and configure per-application how
  the addressbook should work.(Arguable, but possible)
  * Obviate my earlier project [*thirdeye*](https://github.com/eyedeekay/thirdeye)
  by implementing a simpler AddressHelper service(than thirdeye) that can be
  forwarded over i2p and treated like an AddressHelper.(Complete, as much as
  I want it to be anyway, but with an arguable bug where if it's asked for an
  address in it's local addressbook, it will assume that the address is present
  in the asker's local addressbook too. Since the default addressbook is just
  the one that comes with i2pd, this issue may be more-or-less invisible for
  now. See also "Non-Goals")

### Feature Goals

  * Optional DNS-like server.(Not started)
  * Allow the users to configure multiple "master" AddressHelpers, and create
  a way to arrive at a "consensus" based on what they report and longitudinal
  data.(Barely started)
  * Enable automatic **signed** updates using public keys to authenticate
  accounts associated with desired hostnames to create a "Dynamic DNS" type
  setup. This might be useful for services driven by the API, which could change
  destination frequently.(Not started)

## STRATEGIES:

  * Ephemerality(Optional but default). In the default mode of operation, the
  helper doesn't allow the applications using it to cache new addresses. Instead,
  each time it starts, it downloads all the addressbooks it is subscribed to as
  quickly as possible.
  * Rate-Limiting(Optional but default). To confuse timebleed-type attacks, the
  system will use rate-limiting and random delays.
  * Signed Entries(Optional, non-default)
  * Peer-Voting(Optional, non-default)

## Non-Goals:

  * Web interface. It's just a thing you ask for other things. It gives dumb answers
  based on Cheech and Chong skits if you ask it a wrong question. It's for writing
  applications against. Maybe when it's done a web interface can be written
  against it.
