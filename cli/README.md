# ARGUS
This is My Personal Attempt on Argus

This is a PoC - refer to the main branch to see how the actual development looks like.

Should be close to this one anyways.

# TODO for me to be happy with this PoC

- [x] Do a Detailed report showing each requirement, implementation, attestation, and artifacts(logs)
- [x] Add a 'report' command where report types can be specified, eventually files. Let attest command only with summary
- [ ] Implement three use cases using the system as is - see what do I need to change still
- [ ] Make a dependency report where parents status are flagged by the looks of their children status.
- [ ] Add Versioning validation when loading (i.e. old implementations to updated requirements should be invalid - requirement should be also updated from version - refuse to load if that happens)
- [ ] Add a 'show' command that only shows what's loaded in the program

## Running it

```
# Load the state versus the configuration files
./bin/argus load -c ./example/.argus-config.yaml
# Attest resources according to current state
./bin/argus attest -c ./example/.argus-config.yaml
# Reports on the attestation
./bin/argus report -m detailed -o json -c ./example/.argus-config.yaml