### WIP

library capable of loading petri-nets created by http://www.pneditor.org/
adding support for formal validation / simulation

BACKLOG
-------
- [ ] Enforce Roles & Guards in StateMachine transformations
- [ ] add storage interface for persisting state and events
- [ ] add test simulation to check for boundedness in an example file

COMPLETE
--------
- [x] Extend stateMachine to derive a state machine from ptnet.PetriNet
- [x] add better example xml files
- [x] build library for reading & writing pflow petri-nets in xml format
- [x] import pflow as a vector state machine - deal w/ subnets somehow
- [x] add ability to import Conditions/Guards from roles defined in pflow
- [x] figure out what 'static places' are for - learn other pflow details
- [x] add Guards in ptnet.PetriNet 

ICEBOX
------
- [ ] ReferenceArcs in pflow file - do we need them?
- [ ] update previous prototype: FactomProject/ptnet-eventstore to use this lib
- [ ] add golang code generation
