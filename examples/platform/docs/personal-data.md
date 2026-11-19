# Personal data

Personal data is data about a person, such as a name or an address. A service
handles each data class of its tables, and each data class that the services
that it depends on handle:

![](../defs/privacy.of.yaml#dataFromTable)

![](../defs/privacy.of.yaml#dataFromService)

In the past, a customer asked which services had her address. The checkout
service got the address from a different service, and the privacy notice did
not tell the customer. Thus, the privacy notice must name each data class
that a public service handles. Refer to
[publicServiceDeclaresItsData](../defs/privacy.of.yaml#publicServiceDeclaresItsData).
