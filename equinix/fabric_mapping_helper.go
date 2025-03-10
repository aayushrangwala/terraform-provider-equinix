package equinix

import (
	"fmt"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	"log"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func serviceTokenToFabric(serviceTokenRequest []interface{}) (v4.ServiceToken, error) {
	mappedST := v4.ServiceToken{}
	for _, str := range serviceTokenRequest {
		stMap := str.(map[string]interface{})
		stType := stMap["type"].(string)
		uuid := stMap["uuid"].(string)
		if stType != "" {
			if stType != "VC_TOKEN" {
				return v4.ServiceToken{}, fmt.Errorf("invalid service token type in config. Must be: VC_TOKEN; Received: %s", stType)
			}
			stTypeObj := v4.ServiceTokenType(stType)
			mappedST = v4.ServiceToken{Uuid: uuid, Type_: &stTypeObj}
		} else {
			mappedST = v4.ServiceToken{Uuid: uuid}
		}

	}
	return mappedST, nil
}

func additionalInfoTerraToGo(additionalInfoRequest []interface{}) []v4.ConnectionSideAdditionalInfo {
	var mappedaiArray []v4.ConnectionSideAdditionalInfo
	for _, ai := range additionalInfoRequest {
		aiMap := ai.(map[string]interface{})
		key := aiMap["key"].(string)
		value := aiMap["value"].(string)
		mappedai := v4.ConnectionSideAdditionalInfo{Key: key, Value: value}
		mappedaiArray = append(mappedaiArray, mappedai)
	}
	return mappedaiArray
}

func accessPointToFabric(accessPointRequest []interface{}) v4.AccessPoint {
	accessPoint := v4.AccessPoint{}
	for _, ap := range accessPointRequest {
		accessPointMap := ap.(map[string]interface{})
		portList := accessPointMap["port"].(*schema.Set).List()
		profileList := accessPointMap["profile"].(*schema.Set).List()
		locationList := accessPointMap["location"].(*schema.Set).List()
		virtualdeviceList := accessPointMap["virtual_device"].(*schema.Set).List()
		interfaceList := accessPointMap["interface"].(*schema.Set).List()
		networkList := accessPointMap["network"].(*schema.Set).List()
		typeVal := accessPointMap["type"].(string)
		authenticationKey := accessPointMap["authentication_key"].(string)
		if authenticationKey != "" {
			accessPoint.AuthenticationKey = authenticationKey
		}
		providerConnectionId := accessPointMap["provider_connection_id"].(string)
		if providerConnectionId != "" {
			accessPoint.ProviderConnectionId = providerConnectionId
		}
		sellerRegion := accessPointMap["seller_region"].(string)
		if sellerRegion != "" {
			accessPoint.SellerRegion = sellerRegion
		}
		peeringTypeRaw := accessPointMap["peering_type"].(string)
		if peeringTypeRaw != "" {
			peeringType := v4.PeeringType(peeringTypeRaw)
			accessPoint.PeeringType = &peeringType
		}
		cloudRouterRequest := accessPointMap["router"].(*schema.Set).List()
		if len(cloudRouterRequest) == 0 {
			log.Print("[DEBUG] The router attribute was not used, attempting to revert to deprecated gateway attribute")
			cloudRouterRequest = accessPointMap["gateway"].(*schema.Set).List()
		}

		if len(cloudRouterRequest) != 0 {
			mappedGWr := cloudRouterToFabric(cloudRouterRequest)
			if mappedGWr.Uuid != "" {
				accessPoint.Router = &mappedGWr
			}
		}
		apt := v4.AccessPointType(typeVal)
		accessPoint.Type_ = &apt
		if len(portList) != 0 {
			port := portToFabric(portList)
			if port.Uuid != "" {
				accessPoint.Port = &port
			}
		}

		if len(networkList) != 0 {
			network := networkToFabric(networkList)
			if network.Uuid != "" {
				accessPoint.Network = &network
			}
		}
		linkProtocolList := accessPointMap["link_protocol"].(*schema.Set).List()

		if len(linkProtocolList) != 0 {
			slp := linkProtocolToFabric(linkProtocolList)
			if slp.Type_ != nil {
				accessPoint.LinkProtocol = &slp
			}
		}

		if len(profileList) != 0 {
			ssp := simplifiedServiceProfileToFabric(profileList)
			if ssp.Uuid != "" {
				accessPoint.Profile = &ssp
			}
		}

		if len(locationList) != 0 {
			sl := equinix_schema.LocationToFabric(locationList)
			accessPoint.Location = &sl
		}

		if len(virtualdeviceList) != 0 {
			vd := virtualdeviceToFabric(virtualdeviceList)
			accessPoint.VirtualDevice = &vd
		}

		if len(interfaceList) != 0 {
			il := interfaceToFabric(interfaceList)
			accessPoint.Interface_ = &il
		}

	}
	return accessPoint
}

func cloudRouterToFabric(cloudRouterRequest []interface{}) v4.CloudRouter {
	if cloudRouterRequest == nil {
		return v4.CloudRouter{}
	}
	cloudRouterMapped := v4.CloudRouter{}
	for _, crr := range cloudRouterRequest {
		crrMap := crr.(map[string]interface{})
		cruuid := crrMap["uuid"].(string)
		cloudRouterMapped = v4.CloudRouter{Uuid: cruuid}
	}
	return cloudRouterMapped
}

func linkProtocolToFabric(linkProtocolList []interface{}) v4.SimplifiedLinkProtocol {
	slp := v4.SimplifiedLinkProtocol{}
	for _, lp := range linkProtocolList {
		lpMap := lp.(map[string]interface{})
		lpType := lpMap["type"].(string)
		lpVlanSTag := lpMap["vlan_s_tag"].(int)
		lpVlanTag := lpMap["vlan_tag"].(int)
		lpVlanCTag := lpMap["vlan_c_tag"].(int)
		lpt := v4.LinkProtocolType(lpType)
		slp = v4.SimplifiedLinkProtocol{Type_: &lpt, VlanSTag: int32(lpVlanSTag), VlanTag: int32(lpVlanTag), VlanCTag: int32(lpVlanCTag)}
	}
	return slp
}

func networkToFabric(networkList []interface{}) v4.SimplifiedNetwork {
	p := v4.SimplifiedNetwork{}
	for _, pl := range networkList {
		plMap := pl.(map[string]interface{})
		uuid := plMap["uuid"].(string)
		p = v4.SimplifiedNetwork{Uuid: uuid}
	}
	return p
}

func simplifiedServiceProfileToFabric(profileList []interface{}) v4.SimplifiedServiceProfile {
	ssp := v4.SimplifiedServiceProfile{}
	for _, pl := range profileList {
		plMap := pl.(map[string]interface{})
		ptype := plMap["type"].(string)
		spte := v4.ServiceProfileTypeEnum(ptype)
		uuid := plMap["uuid"].(string)
		ssp = v4.SimplifiedServiceProfile{Uuid: uuid, Type_: &spte}

	}
	return ssp
}

func virtualdeviceToFabric(virtualdeviceList []interface{}) v4.VirtualDevice {
	vd := v4.VirtualDevice{}
	for _, ll := range virtualdeviceList {
		llMap := ll.(map[string]interface{})
		hr := llMap["href"].(string)
		tp := llMap["type"].(string)
		ud := llMap["uuid"].(string)
		na := llMap["name"].(string)
		vd = v4.VirtualDevice{Href: hr, Type_: tp, Uuid: ud, Name: na}
	}
	return vd
}

func interfaceToFabric(interfaceList []interface{}) v4.ModelInterface {
	il := v4.ModelInterface{}
	for _, ll := range interfaceList {
		llMap := ll.(map[string]interface{})
		ud := llMap["uuid"].(string)
		tp := llMap["type"].(string)
		id := llMap["id"].(int)
		il = v4.ModelInterface{Type_: tp, Uuid: ud, Id: int32(id)}
	}
	return il
}

func operationToTerra(operation *v4.ConnectionOperation) *schema.Set {
	if operation == nil {
		return nil
	}
	operations := []*v4.ConnectionOperation{operation}
	mappedOperations := make([]interface{}, len(operations))
	for _, operation := range operations {
		mappedOperation := make(map[string]interface{})
		mappedOperation["provider_status"] = string(*operation.ProviderStatus)
		mappedOperation["equinix_status"] = string(*operation.EquinixStatus)
		if operation.Errors != nil {
			mappedOperation["errors"] = equinix_schema.ErrorToTerra(operation.Errors)
		}
		mappedOperations = append(mappedOperations, mappedOperation)
	}
	operationSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: operationSch()}),
		mappedOperations,
	)
	return operationSet
}

func serviceTokenToTerra(serviceToken *v4.ServiceToken) *schema.Set {
	if serviceToken == nil {
		return nil
	}
	serviceTokens := []*v4.ServiceToken{serviceToken}
	mappedServiceTokens := make([]interface{}, len(serviceTokens))
	for _, serviceToken := range serviceTokens {
		mappedServiceToken := make(map[string]interface{})
		if serviceToken.Type_ != nil {
			mappedServiceToken["type"] = string(*serviceToken.Type_)
		}
		mappedServiceToken["href"] = serviceToken.Href
		mappedServiceToken["uuid"] = serviceToken.Uuid
		mappedServiceTokens = append(mappedServiceTokens, mappedServiceToken)
	}
	serviceTokenSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: serviceTokenSch()}),
		mappedServiceTokens,
	)
	return serviceTokenSet
}

func connectionSideToTerra(connectionSide *v4.ConnectionSide) *schema.Set {
	connectionSides := []*v4.ConnectionSide{connectionSide}
	mappedConnectionSides := make([]interface{}, len(connectionSides))
	for _, connectionSide := range connectionSides {
		mappedConnectionSide := make(map[string]interface{})
		serviceTokenSet := serviceTokenToTerra(connectionSide.ServiceToken)
		if serviceTokenSet != nil {
			mappedConnectionSide["service_token"] = serviceTokenSet
		}
		mappedConnectionSide["access_point"] = accessPointToTerra(connectionSide.AccessPoint)
		mappedConnectionSides = append(mappedConnectionSides, mappedConnectionSide)
	}
	connectionSideSet := schema.NewSet(
		schema.HashResource(connectionSideSch()),
		mappedConnectionSides,
	)
	return connectionSideSet
}

func additionalInfoToTerra(additionalInfol []v4.ConnectionSideAdditionalInfo) []map[string]interface{} {
	if additionalInfol == nil {
		return nil
	}
	mappedadditionalInfol := make([]map[string]interface{}, len(additionalInfol))
	for index, additionalInfo := range additionalInfol {
		mappedadditionalInfol[index] = map[string]interface{}{
			"key":   additionalInfo.Key,
			"value": additionalInfo.Value,
		}
	}
	return mappedadditionalInfol
}

func cloudRouterToTerra(cloudRouter *v4.CloudRouter) *schema.Set {
	if cloudRouter == nil {
		return nil
	}
	cloudRouters := []*v4.CloudRouter{cloudRouter}
	mappedCloudRouters := make([]interface{}, len(cloudRouters))
	for _, cloudRouter := range cloudRouters {
		mappedCloudRouter := make(map[string]interface{})
		mappedCloudRouter["uuid"] = cloudRouter.Uuid
		mappedCloudRouter["href"] = cloudRouter.Href
		mappedCloudRouters = append(mappedCloudRouters, mappedCloudRouter)
	}
	linkedProtocolSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: equinix_schema.ProjectSch()}),
		mappedCloudRouters)
	return linkedProtocolSet
}

func virtualDeviceToTerra(virtualDevice *v4.VirtualDevice) *schema.Set {
	if virtualDevice == nil {
		return nil
	}
	virtualDevices := []*v4.VirtualDevice{virtualDevice}
	mappedVirtualDevices := make([]interface{}, len(virtualDevices))
	for _, virtualDevice := range virtualDevices {
		mappedVirtualDevice := make(map[string]interface{})
		mappedVirtualDevice["name"] = virtualDevice.Name
		mappedVirtualDevice["href"] = virtualDevice.Href
		mappedVirtualDevice["type"] = virtualDevice.Type_
		mappedVirtualDevice["uuid"] = virtualDevice.Uuid
		mappedVirtualDevices = append(mappedVirtualDevices, mappedVirtualDevice)
	}
	virtualDeviceSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: accessPointVirtualDeviceSch()}),
		mappedVirtualDevices)
	return virtualDeviceSet
}

func interfaceToTerra(mInterface *v4.ModelInterface) *schema.Set {
	if mInterface == nil {
		return nil
	}
	mInterfaces := []*v4.ModelInterface{mInterface}
	mappedMInterfaces := make([]interface{}, len(mInterfaces))
	for _, mInterface := range mInterfaces {
		mappedMInterface := make(map[string]interface{})
		mappedMInterface["id"] = int(mInterface.Id)
		mappedMInterface["type"] = mInterface.Type_
		mappedMInterface["uuid"] = mInterface.Uuid
		mappedMInterfaces = append(mappedMInterfaces, mappedMInterface)
	}
	mInterfaceSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: accessPointInterface()}),
		mappedMInterfaces)
	return mInterfaceSet
}

func accessPointToTerra(accessPoint *v4.AccessPoint) *schema.Set {
	accessPoints := []*v4.AccessPoint{accessPoint}
	mappedAccessPoints := make([]interface{}, len(accessPoints))
	for _, accessPoint := range accessPoints {
		mappedAccessPoint := make(map[string]interface{})
		if accessPoint.Type_ != nil {
			mappedAccessPoint["type"] = string(*accessPoint.Type_)
		}
		if accessPoint.Account != nil {
			mappedAccessPoint["account"] = equinix_schema.AccountToTerra(accessPoint.Account)
		}
		if accessPoint.Location != nil {
			mappedAccessPoint["location"] = equinix_schema.LocationToTerra(accessPoint.Location)
		}
		if accessPoint.Port != nil {
			mappedAccessPoint["port"] = portToTerra(accessPoint.Port)
		}
		if accessPoint.Profile != nil {
			mappedAccessPoint["profile"] = simplifiedServiceProfileToTerra(accessPoint.Profile)
		}
		if accessPoint.Router != nil {
			mappedAccessPoint["router"] = cloudRouterToTerra(accessPoint.Router)
			mappedAccessPoint["gateway"] = cloudRouterToTerra(accessPoint.Router)
		}
		if accessPoint.LinkProtocol != nil {
			mappedAccessPoint["link_protocol"] = linkedProtocolToTerra(*accessPoint.LinkProtocol)
		}
		if accessPoint.VirtualDevice != nil {
			mappedAccessPoint["virtual_device"] = virtualDeviceToTerra(accessPoint.VirtualDevice)
		}
		if accessPoint.Interface_ != nil {
			mappedAccessPoint["interface"] = interfaceToTerra(accessPoint.Interface_)
		}
		mappedAccessPoint["seller_region"] = accessPoint.SellerRegion
		if accessPoint.PeeringType != nil {
			mappedAccessPoint["peering_type"] = string(*accessPoint.PeeringType)
		}
		mappedAccessPoint["authentication_key"] = accessPoint.AuthenticationKey
		mappedAccessPoint["provider_connection_id"] = accessPoint.ProviderConnectionId
		mappedAccessPoints = append(mappedAccessPoints, mappedAccessPoint)
	}
	accessPointSet := schema.NewSet(
		schema.HashResource(accessPointSch()),
		mappedAccessPoints,
	)
	return accessPointSet
}

func linkedProtocolToTerra(linkedProtocol v4.SimplifiedLinkProtocol) *schema.Set {
	linkedProtocols := []v4.SimplifiedLinkProtocol{linkedProtocol}
	mappedLinkedProtocols := make([]interface{}, len(linkedProtocols))
	for _, linkedProtocol := range linkedProtocols {
		mappedLinkedProtocol := make(map[string]interface{})
		mappedLinkedProtocol["type"] = string(*linkedProtocol.Type_)
		mappedLinkedProtocol["vlan_tag"] = int(linkedProtocol.VlanTag)
		mappedLinkedProtocol["vlan_s_tag"] = int(linkedProtocol.VlanSTag)
		mappedLinkedProtocol["vlan_c_tag"] = int(linkedProtocol.VlanCTag)
		mappedLinkedProtocols = append(mappedLinkedProtocols, mappedLinkedProtocol)
	}
	linkedProtocolSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: accessPointLinkProtocolSch()}),
		mappedLinkedProtocols)
	return linkedProtocolSet
}

func simplifiedServiceProfileToTerra(profile *v4.SimplifiedServiceProfile) *schema.Set {
	profiles := []*v4.SimplifiedServiceProfile{profile}
	mappedProfiles := make([]interface{}, len(profiles))
	for _, profile := range profiles {
		mappedProfile := make(map[string]interface{})
		mappedProfile["href"] = profile.Href
		mappedProfile["type"] = string(*profile.Type_)
		mappedProfile["name"] = profile.Name
		mappedProfile["uuid"] = profile.Uuid
		mappedProfile["access_point_type_configs"] = accessPointTypeConfigToTerra(profile.AccessPointTypeConfigs)
		mappedProfiles = append(mappedProfiles, mappedProfile)
	}

	profileSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: serviceProfileSch()}),
		mappedProfiles,
	)
	return profileSet
}

func accessPointTypeConfigToTerra(spAccessPointTypes []v4.ServiceProfileAccessPointType) []interface{} {
	mappedSpAccessPointTypes := make([]interface{}, len(spAccessPointTypes))
	for index, spAccessPointType := range spAccessPointTypes {
		mappedSpAccessPointTypes[index] = map[string]interface{}{
			"type":                             string(*spAccessPointType.Type_),
			"uuid":                             spAccessPointType.Uuid,
			"allow_remote_connections":         spAccessPointType.AllowRemoteConnections,
			"allow_custom_bandwidth":           spAccessPointType.AllowCustomBandwidth,
			"allow_bandwidth_auto_approval":    spAccessPointType.AllowBandwidthAutoApproval,
			"enable_auto_generate_service_key": spAccessPointType.EnableAutoGenerateServiceKey,
			"connection_redundancy_required":   spAccessPointType.ConnectionRedundancyRequired,
			"connection_label":                 spAccessPointType.ConnectionLabel,
			"api_config":                       apiConfigToTerra(spAccessPointType.ApiConfig),
			"authentication_key":               authenticationKeyToTerra(spAccessPointType.AuthenticationKey),
			"supported_bandwidths":             supportedBandwidthsToTerra(spAccessPointType.SupportedBandwidths),
		}
	}

	return mappedSpAccessPointTypes
}

func apiConfigToTerra(apiConfig *v4.ApiConfig) *schema.Set {
	apiConfigs := []*v4.ApiConfig{apiConfig}
	mappedApiConfigs := make([]interface{}, len(apiConfigs))
	for _, apiConfig := range apiConfigs {
		mappedApiConfig := make(map[string]interface{})
		mappedApiConfig["api_available"] = apiConfig.ApiAvailable
		mappedApiConfig["equinix_managed_vlan"] = apiConfig.EquinixManagedVlan
		mappedApiConfig["bandwidth_from_api"] = apiConfig.BandwidthFromApi
		mappedApiConfig["integration_id"] = apiConfig.IntegrationId
		mappedApiConfig["equinix_managed_port"] = apiConfig.EquinixManagedPort
		mappedApiConfigs = append(mappedApiConfigs, mappedApiConfig)
	}
	apiConfigSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createApiConfigSch()}),
		mappedApiConfigs)
	return apiConfigSet
}

func authenticationKeyToTerra(authenticationKey *v4.AuthenticationKey) *schema.Set {
	authenticationKeys := []*v4.AuthenticationKey{authenticationKey}
	mappedAuthenticationKeys := make([]interface{}, len(authenticationKeys))
	for _, authenticationKey := range authenticationKeys {
		mappedAuthenticationKey := make(map[string]interface{})
		mappedAuthenticationKey["required"] = authenticationKey.Required
		mappedAuthenticationKey["label"] = authenticationKey.Label
		mappedAuthenticationKey["description"] = authenticationKey.Description
		mappedAuthenticationKeys = append(mappedAuthenticationKeys, mappedAuthenticationKey)
	}
	apiConfigSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createAuthenticationKeySch()}),
		mappedAuthenticationKeys)
	return apiConfigSet
}

func supportedBandwidthsToTerra(supportedBandwidths *[]int32) []interface{} {
	if supportedBandwidths == nil {
		return nil
	}
	mappedSupportedBandwidths := make([]interface{}, len(*supportedBandwidths))
	for _, bandwidth := range *supportedBandwidths {
		mappedSupportedBandwidths = append(mappedSupportedBandwidths, int(bandwidth))
	}
	return mappedSupportedBandwidths
}

func routingProtocolDirectIpv4ToFabric(routingProtocolDirectIpv4Request []interface{}) v4.DirectConnectionIpv4 {
	mappedRpDirectIpv4 := v4.DirectConnectionIpv4{}
	for _, str := range routingProtocolDirectIpv4Request {
		directIpv4Map := str.(map[string]interface{})
		equinixIfaceIp := directIpv4Map["equinix_iface_ip"].(string)

		mappedRpDirectIpv4 = v4.DirectConnectionIpv4{EquinixIfaceIp: equinixIfaceIp}
	}
	return mappedRpDirectIpv4
}

func routingProtocolDirectIpv6ToFabric(routingProtocolDirectIpv6Request []interface{}) v4.DirectConnectionIpv6 {
	mappedRpDirectIpv6 := v4.DirectConnectionIpv6{}
	for _, str := range routingProtocolDirectIpv6Request {
		directIpv6Map := str.(map[string]interface{})
		equinixIfaceIp := directIpv6Map["equinix_iface_ip"].(string)

		mappedRpDirectIpv6 = v4.DirectConnectionIpv6{EquinixIfaceIp: equinixIfaceIp}
	}
	return mappedRpDirectIpv6
}

func routingProtocolBgpIpv4ToFabric(routingProtocolBgpIpv4Request []interface{}) v4.BgpConnectionIpv4 {
	mappedRpBgpIpv4 := v4.BgpConnectionIpv4{}
	for _, str := range routingProtocolBgpIpv4Request {
		bgpIpv4Map := str.(map[string]interface{})
		customerPeerIp := bgpIpv4Map["customer_peer_ip"].(string)
		enabled := bgpIpv4Map["enabled"].(bool)

		mappedRpBgpIpv4 = v4.BgpConnectionIpv4{CustomerPeerIp: customerPeerIp, Enabled: enabled}
	}
	return mappedRpBgpIpv4
}

func routingProtocolBgpIpv6ToFabric(routingProtocolBgpIpv6Request []interface{}) v4.BgpConnectionIpv6 {
	mappedRpBgpIpv6 := v4.BgpConnectionIpv6{}
	for _, str := range routingProtocolBgpIpv6Request {
		bgpIpv6Map := str.(map[string]interface{})
		customerPeerIp := bgpIpv6Map["customer_peer_ip"].(string)
		enabled := bgpIpv6Map["enabled"].(bool)

		mappedRpBgpIpv6 = v4.BgpConnectionIpv6{CustomerPeerIp: customerPeerIp, Enabled: enabled}
	}
	return mappedRpBgpIpv6
}

func routingProtocolBfdToFabric(routingProtocolBfdRequest []interface{}) v4.RoutingProtocolBfd {
	mappedRpBfd := v4.RoutingProtocolBfd{}
	for _, str := range routingProtocolBfdRequest {
		rpBfdMap := str.(map[string]interface{})
		bfdEnabled := rpBfdMap["enabled"].(bool)
		bfdInterval := rpBfdMap["interval"].(string)

		mappedRpBfd = v4.RoutingProtocolBfd{Enabled: bfdEnabled, Interval: bfdInterval}
	}
	return mappedRpBfd
}

func routingProtocolChangeToFabric(routingProtocolChangeRequest []interface{}) v4.RoutingProtocolChange {
	mappedRpChange := v4.RoutingProtocolChange{}
	for _, str := range routingProtocolChangeRequest {
		rpChangeMap := str.(map[string]interface{})
		uuid := rpChangeMap["uuid"].(string)
		rpChangeType := rpChangeMap["type"].(string)

		mappedRpChange = v4.RoutingProtocolChange{Uuid: uuid, Type_: rpChangeType}
	}
	return mappedRpChange
}

func routingProtocolDirectTypeToTerra(routingProtocolDirect *v4.RoutingProtocolDirectType) *schema.Set {
	if routingProtocolDirect == nil {
		return nil
	}
	routingProtocolDirects := []*v4.RoutingProtocolDirectType{routingProtocolDirect}
	mappedDirects := make([]interface{}, len(routingProtocolDirects))
	for _, routingProtocolDirect := range routingProtocolDirects {
		mappedDirect := make(map[string]interface{})
		mappedDirect["type"] = routingProtocolDirect.Type_
		mappedDirect["name"] = routingProtocolDirect.Name
		if routingProtocolDirect.DirectIpv4 != nil {
			mappedDirect["direct_ipv4"] = routingProtocolDirectConnectionIpv4ToTerra(routingProtocolDirect.DirectIpv4)
		}
		if routingProtocolDirect.DirectIpv6 != nil {
			mappedDirect["direct_ipv6"] = routingProtocolDirectConnectionIpv6ToTerra(routingProtocolDirect.DirectIpv6)
		}
		mappedDirects = append(mappedDirects, mappedDirect)
	}
	rpDirectSet := schema.NewSet(
		schema.HashResource(createRoutingProtocolDirectTypeRes),
		mappedDirects,
	)

	return rpDirectSet
}

func routingProtocolDirectConnectionIpv4ToTerra(routingProtocolDirectIpv4 *v4.DirectConnectionIpv4) *schema.Set {
	if routingProtocolDirectIpv4 == nil {
		return nil
	}
	routingProtocolDirectIpv4s := []*v4.DirectConnectionIpv4{routingProtocolDirectIpv4}
	mappedDirectIpv4s := make([]interface{}, len(routingProtocolDirectIpv4s))
	for i, routingProtocolDirectIpv4 := range routingProtocolDirectIpv4s {
		mappedDirectIpv4s[i] = map[string]interface{}{
			"equinix_iface_ip": routingProtocolDirectIpv4.EquinixIfaceIp,
		}
	}
	rpDirectIpv4Set := schema.NewSet(
		schema.HashResource(createDirectConnectionIpv4Res),
		mappedDirectIpv4s,
	)
	return rpDirectIpv4Set
}

func routingProtocolDirectConnectionIpv6ToTerra(routingProtocolDirectIpv6 *v4.DirectConnectionIpv6) *schema.Set {
	if routingProtocolDirectIpv6 == nil {
		return nil
	}
	routingProtocolDirectIpv6s := []*v4.DirectConnectionIpv6{routingProtocolDirectIpv6}
	mappedDirectIpv6s := make([]interface{}, len(routingProtocolDirectIpv6s))
	for i, routingProtocolDirectIpv6 := range routingProtocolDirectIpv6s {
		mappedDirectIpv6s[i] = map[string]interface{}{
			"equinix_iface_ip": routingProtocolDirectIpv6.EquinixIfaceIp,
		}
	}
	rpDirectIpv6Set := schema.NewSet(
		schema.HashResource(createDirectConnectionIpv6Res),
		mappedDirectIpv6s,
	)
	return rpDirectIpv6Set
}

func routingProtocolBgpTypeToTerra(routingProtocolBgp *v4.RoutingProtocolBgpType) *schema.Set {
	if routingProtocolBgp == nil {
		return nil
	}
	routingProtocolBgps := []*v4.RoutingProtocolBgpType{routingProtocolBgp}
	mappedBgps := make([]interface{}, len(routingProtocolBgps))
	for _, routingProtocolBgp := range routingProtocolBgps {
		mappedBgp := make(map[string]interface{})
		mappedBgp["type"] = routingProtocolBgp.Type_
		mappedBgp["name"] = routingProtocolBgp.Name
		if routingProtocolBgp.BgpIpv4 != nil {
			mappedBgp["bgp_ipv4"] = routingProtocolBgpConnectionIpv4ToTerra(routingProtocolBgp.BgpIpv4)
		}
		if routingProtocolBgp.BgpIpv6 != nil {
			mappedBgp["bgp_ipv6"] = routingProtocolBgpConnectionIpv6ToTerra(routingProtocolBgp.BgpIpv6)
		}
		mappedBgp["customer_asn"] = routingProtocolBgp.CustomerAsn
		mappedBgp["bgp_auth_key"] = routingProtocolBgp.BgpAuthKey
		if routingProtocolBgp.Bfd != nil {
			mappedBgp["bfd"] = routingProtocolBfdToTerra(routingProtocolBgp.Bfd)
		}

		mappedBgps = append(mappedBgps, mappedBgp)

	}
	rpBgpSet := schema.NewSet(
		schema.HashResource(createRoutingProtocolBgpTypeRes),
		mappedBgps,
	)

	return rpBgpSet
}

func routingProtocolBgpConnectionIpv4ToTerra(routingProtocolBgpIpv4 *v4.BgpConnectionIpv4) *schema.Set {
	if routingProtocolBgpIpv4 == nil {
		return nil
	}
	routingProtocolBgpIpv4s := []*v4.BgpConnectionIpv4{routingProtocolBgpIpv4}
	mappedBgpIpv4s := make([]interface{}, len(routingProtocolBgpIpv4s))
	for i, routingProtocolBgpIpv4 := range routingProtocolBgpIpv4s {
		mappedBgpIpv4s[i] = map[string]interface{}{
			"customer_peer_ip": routingProtocolBgpIpv4.CustomerPeerIp,
			"equinix_peer_ip":  routingProtocolBgpIpv4.EquinixPeerIp,
			"enabled":          routingProtocolBgpIpv4.Enabled,
		}
	}
	rpBgpIpv4Set := schema.NewSet(
		schema.HashResource(createBgpConnectionIpv4Res),
		mappedBgpIpv4s,
	)
	return rpBgpIpv4Set
}

func routingProtocolBgpConnectionIpv6ToTerra(routingProtocolBgpIpv6 *v4.BgpConnectionIpv6) *schema.Set {
	if routingProtocolBgpIpv6 == nil {
		return nil
	}
	routingProtocolBgpIpv6s := []*v4.BgpConnectionIpv6{routingProtocolBgpIpv6}
	mappedBgpIpv6s := make([]interface{}, len(routingProtocolBgpIpv6s))
	for i, routingProtocolBgpIpv6 := range routingProtocolBgpIpv6s {
		mappedBgpIpv6s[i] = map[string]interface{}{
			"customer_peer_ip": routingProtocolBgpIpv6.CustomerPeerIp,
			"equinix_peer_ip":  routingProtocolBgpIpv6.EquinixPeerIp,
			"enabled":          routingProtocolBgpIpv6.Enabled,
		}
	}
	rpBgpIpv6Set := schema.NewSet(
		schema.HashResource(createBgpConnectionIpv6Res),
		mappedBgpIpv6s,
	)
	return rpBgpIpv6Set
}

func routingProtocolBfdToTerra(routingProtocolBfd *v4.RoutingProtocolBfd) *schema.Set {
	if routingProtocolBfd == nil {
		return nil
	}
	routingProtocolBfds := []*v4.RoutingProtocolBfd{routingProtocolBfd}
	mappedRpBfds := make([]interface{}, len(routingProtocolBfds))
	for i, routingProtocolBfd := range routingProtocolBfds {
		mappedRpBfds[i] = map[string]interface{}{
			"enabled":  routingProtocolBfd.Enabled,
			"interval": routingProtocolBfd.Interval,
		}
	}
	rpBfdSet := schema.NewSet(
		schema.HashResource(createRoutingProtocolBfdRes),
		mappedRpBfds,
	)
	return rpBfdSet
}

func routingProtocolOperationToTerra(routingProtocolOperation *v4.RoutingProtocolOperation) *schema.Set {
	if routingProtocolOperation == nil {
		return nil
	}
	routingProtocolOperations := []*v4.RoutingProtocolOperation{routingProtocolOperation}
	mappedRpOperations := make([]interface{}, len(routingProtocolOperations))
	for _, routingProtocolOperation := range routingProtocolOperations {
		mappedRpOperation := make(map[string]interface{})
		if routingProtocolOperation.Errors != nil {
			mappedRpOperation["errors"] = equinix_schema.ErrorToTerra(routingProtocolOperation.Errors)
		}
		mappedRpOperations = append(mappedRpOperations, mappedRpOperation)
	}
	rpOperationSet := schema.NewSet(
		schema.HashResource(createRoutingProtocolOperationRes),
		mappedRpOperations,
	)
	return rpOperationSet
}

func routingProtocolChangeToTerra(routingProtocolChange *v4.RoutingProtocolChange) *schema.Set {
	if routingProtocolChange == nil {
		return nil
	}
	routingProtocolChanges := []*v4.RoutingProtocolChange{routingProtocolChange}
	mappedRpChanges := make([]interface{}, len(routingProtocolChanges))
	for i, rpChanges := range routingProtocolChanges {
		mappedRpChanges[i] = map[string]interface{}{
			"uuid": rpChanges.Uuid,
			"type": rpChanges.Type_,
			"href": rpChanges.Href,
		}
	}
	rpChangeSet := schema.NewSet(
		schema.HashResource(createRoutingProtocolChangeRes),
		mappedRpChanges,
	)
	return rpChangeSet
}

func getRoutingProtocolPatchUpdateRequest(rp v4.RoutingProtocolData, d *schema.ResourceData) (v4.ConnectionChangeOperation, error) {
	changeOps := v4.ConnectionChangeOperation{}
	existingBgpIpv4Status := rp.BgpIpv4.Enabled
	existingBgpIpv6Status := rp.BgpIpv6.Enabled
	updateBgpIpv4Status := d.Get("rp.BgpIpv4.Enabled")
	updateBgpIpv6Status := d.Get("rp.BgpIpv6.Enabled")

	log.Printf("existing BGP IPv4 Status: %t, existing BGP IPv6 Status: %t, Update BGP IPv4 Request: %t, Update BGP Ipv6 Request: %t",
		existingBgpIpv4Status, existingBgpIpv6Status, updateBgpIpv4Status, updateBgpIpv6Status)

	if existingBgpIpv4Status != updateBgpIpv4Status {
		changeOps = v4.ConnectionChangeOperation{Op: "replace", Path: "/bgpIpv4/enabled", Value: updateBgpIpv4Status}
	} else if existingBgpIpv6Status != updateBgpIpv6Status {
		changeOps = v4.ConnectionChangeOperation{Op: "replace", Path: "/bgpIpv6/enabled", Value: updateBgpIpv6Status}
	} else {
		return changeOps, fmt.Errorf("nothing to update for the routing protocol %s", rp.RoutingProtocolBgpData.Uuid)
	}
	return changeOps, nil
}

func getUpdateRequests(conn v4.Connection, d *schema.ResourceData) ([][]v4.ConnectionChangeOperation, error) {
	var changeOps [][]v4.ConnectionChangeOperation
	existingName := conn.Name
	existingBandwidth := int(conn.Bandwidth)
	updateNameVal := d.Get("name").(string)
	updateBandwidthVal := d.Get("bandwidth").(int)
	additionalInfo := d.Get("additional_info").([]interface{})

	awsSecrets, hasAWSSecrets := additionalInfoContainsAWSSecrets(additionalInfo)

	if existingName != updateNameVal {
		changeOps = append(changeOps, []v4.ConnectionChangeOperation{
			{
				Op:    "replace",
				Path:  "/name",
				Value: updateNameVal,
			},
		})
	}

	if existingBandwidth != updateBandwidthVal {
		changeOps = append(changeOps, []v4.ConnectionChangeOperation{
			{
				Op:    "replace",
				Path:  "/bandwidth",
				Value: updateBandwidthVal,
			},
		})
	}

	if *conn.Operation.ProviderStatus == v4.PENDING_APPROVAL_ProviderStatus && hasAWSSecrets {
		changeOps = append(changeOps, []v4.ConnectionChangeOperation{
			{
				Op:    "add",
				Path:  "",
				Value: map[string]interface{}{"additionalInfo": awsSecrets},
			},
		})
	}

	if len(changeOps) == 0 {
		return changeOps, fmt.Errorf("nothing to update for the connection %s", existingName)
	}

	return changeOps, nil
}
