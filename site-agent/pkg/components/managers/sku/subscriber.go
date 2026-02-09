package sku

// RegisterSubscriber registers the SKU workflows/activities with the Temporal client
// Note: SKU does not have create/update/delete capabilities, so no subscriber workflows are registered
func (api *API) RegisterSubscriber() error {
	// Register the subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("SKU: Registering the subscribers")
	// Note: SKU is read-only, no CRUD workflows to register
	ManagerAccess.Data.EB.Log.Info().Msg("SKU: No CRUD workflows for SKU (read-only resource)")

	return nil
}
