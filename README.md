# Exhaustive-enum

Exhaustive-enum is an exhaustive enum checker for Go. It can be used to make sure you don't miss any cases when switching over an enumeration.

## Method

Find switch statements. For each:

* Find type of expression being switched over.
* Verify if the type is an enum:
  1. `exhaustive-enum` annotation.
  2. Type has to be a value type (not struct/interface).
  3. Find all exported values in the package of the type definition.
* Verify that all the enum cases are in the switch, or that there's a default case.
