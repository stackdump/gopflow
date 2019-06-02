### WIP

library capable of loading petri-nets created by http://www.pneditor.org/
adding support for formal validation / simulation

BACKLOG
-------
- [ ] update Guards in ptnet.PetriNet to be stored as a map
- [ ] ReferenceArcs in pflow file - do we need them?
- [ ] add test simulation to check for boundedness in an example file
- [ ] Extend stateMachine to derive a state machine from ptnet.PetriNet
- [ ] update previous prototype: FactomProject/ptnet-eventstore to use this lib

COMPLETE
--------
- [x] add better example xml files
- [x] build library for reading & writing pflow petri-nets in xml format
- [x] import pflow as a vector state machine - deal w/ subnets somehow
- [x] add ability to import Conditions/Guards from roles defined in pflow
- [x] figure out what 'static places' are for - learn other pflow details

ICEBOX
------
- [ ] add golang code generation
