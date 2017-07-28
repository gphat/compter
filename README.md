# compter

Pay Your Debts!

# What Is It?

Compter is a web service that tracks a "debt", which represents a thing that we want to complete. It is a common observability problem to need the status of a promised action. For example: did that cron job run, has that recurring event happened a least one in the last hour, &c!

# API

All the API endpoints return JSON.

The debt entity looks like this:
```
type Debt struct {
	ID       int64     `json:"id"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Payments int64     `json:"payments"`
}
```

## PUT /debts

Creates a debt, returning an ID.

## GET /debts/:id

Returns the debt, if present, for the given ID.

## POST /debts/:id

Updates the debt's `Updated` time and increases it's `Payments` by 1.

## GET /debts

Lists all debts.

# TODO

* Metrics emission: If Compter periodically emits metrics you can alert on "has this happened?". Also, if emitted as a monotonic counter a rate could be used, as the rate of payment should equal the expected period e.g. once per hour, etc.
* Event emission: Some systems may prefer an analytics approach, emitting an event with all the attributes.
* Schedules/Interval: How often should a "payment" occur? (Could simplify to a deadline timestamp)
* Persistence: since this PoC is only in memory. Using a backend like Pg or Redis would allow persistent and more scaling.
* Multi-stage payments with timeout: A debt may require multiple payments, imagine an event like a deploy that starts then moves through a series of steps. Each "payment" would have a name (i.e. `fetched_info`, `copied_artifact`, `restarted_service`) and a final payoff. If the timeout is reached and no payoff has occurred, such information is emitted and the debt cleared. All payments are made via a `POST` with the debt's ID
* Vivifying debts: If a POST is made to a non-existent debt, create one and differentiate from an update via status code. This allows recovery from lost debts, as it is assumed that the system normally is working well. (e.g a storage failure loses some/all debts, but cron *normally* works so they get reentered). This would require that POSTs include all data needed to recreate the debt, mainly the interval.

# Name

A [compter](https://en.wikipedia.org/wiki/Compter) was a type of prison for debtors, mostly [during the 18th and 19th centuries](https://en.wikipedia.org/wiki/Debtors%27_prison#Great_Britain_.28later_the_United_Kingdom.29). Those in debt remained in the compter until they could cover their debts.
