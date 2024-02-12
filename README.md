# MOOgraph :milk: :mag: :cow:

A tool for looking at the contents of moo databases

Very WIP - not even full POC at this point.

TODOs
- Backfill a few tests
  - Start with some acceptance-style tests using the test DB you found on github
- Run thru object parsing code, make sure youre parsing all the things, data is valid
- Implement proper list structures for object definitions
- parallelize object parsing
- parallelize above-mentioned bounds checking? waste of time for now im sure
- ???

DONE
- Find bounds of all objects in object block (wut? not sure what i meant here but i done did it)
- Implement sequential search for object definition bounds - eliminate global state
