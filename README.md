API abstraction layer into useful information from vCenter.

Currently only supports the /haReport endpoint which will give a list of
any VMs have that had an HA event in the last 12 hours (eventually this
will be configurable).

All vCenters listed in the _vcenter._tcp service list are queried.
