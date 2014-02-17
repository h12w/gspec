Design of gspec/expectation
===========================

Goals:
* Provide structured and helpful error messages.
* Extensible.
* Can be used alone with "go test".
* Fluent interface.
* Do not panic.

Error Message
-------------
An error message should tell the developer clearly:
* Where the failure is
* What is expected
* What is actually got

Should a customizable description supported? No!
* In a BDD-style framework, such descriptions already included in the "it" test
  group. e.g. it("should xxx", func() { expect ... }). So providing such
  optional arguments will be a duplication.
* When used directly with "go test", such messages can be easily written as a
  comment above the expectation line of code, and printed out when the test
  fails.

An error message should be provided in the form of an error object.
* Plain string message can be easily obtained by method Error().
* Detail structure can still be hold within the object, so that an aggresive 
  reporter can use such information to provide an enhanced visualization such
  as colorful highlighting.

List of expectations
--------------------
This is a complete list of expectations, "+" means already supported:
* +To: general expectation
* +Panic: panic and panic with an object
* +Equal: deep equal
*  Is: shallow equal
*  Order comparision: >, <, >=, <=, Within(delta).Of(value)
*  Composition: Not, And, Or
*  True/False: can be supported by Is(true), Is(false)
*  Nil: can be supported by Is(nil)
*  Empty: empty for container types
*  Error: can be supported by Equal
*  Implements: can be supported by IsType
*  IsType: if a interface is of a type
*  String: Contains, Match...
*  Collection: HasLen, All, Any, Exact, InOrder, InPartialOrder ...

