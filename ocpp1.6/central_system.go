package ocpp16

import (
	"fmt"
	"reflect"

	"github.com/lorenzodonini/ocpp-go/internal/callbackqueue"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/certificates"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/extendedtriggermessage"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/logging"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/securefirmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/security"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
)

type centralSystem struct {
	server                *ocppj.Server
	coreHandler           core.CentralSystemHandler
	localAuthListHandler  localauth.CentralSystemHandler
	firmwareHandler       firmware.CentralSystemHandler
	reservationHandler    reservation.CentralSystemHandler
	remoteTriggerHandler  remotetrigger.CentralSystemHandler
	smartChargingHandler  smartcharging.CentralSystemHandler
	logHandler            logging.CentralSystemHandler
	securityHandler       security.CentralSystemHandler
	secureFirmwareHandler securefirmware.CentralSystemHandler
	callbackQueue         callbackqueue.CallbackQueue
	errC                  chan error
}

func newCentralSystem(server *ocppj.Server) centralSystem {
	if server == nil {
		panic("server must not be nil")
	}
	server.SetDialect(ocpp.V16)
	return centralSystem{
		server:        server,
		callbackQueue: callbackqueue.New(),
	}
}

func (cs *centralSystem) error(err error) {
	if cs.errC != nil {
		cs.errC <- err
	}
}

func (cs *centralSystem) Errors() <-chan error {
	if cs.errC == nil {
		cs.errC = make(chan error, 1)
	}
	return cs.errC
}

func (cs *centralSystem) ChangeAvailability(clientId string, callback func(confirmation *core.ChangeAvailabilityConfirmation, err error), connectorId int, availabilityType core.AvailabilityType, props ...func(request *core.ChangeAvailabilityRequest)) error {
	request := core.NewChangeAvailabilityRequest(connectorId, availabilityType)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*core.ChangeAvailabilityConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) ChangeConfiguration(clientId string, callback func(confirmation *core.ChangeConfigurationConfirmation, err error), key string, value string, props ...func(request *core.ChangeConfigurationRequest)) error {
	request := core.NewChangeConfigurationRequest(key, value)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*core.ChangeConfigurationConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) ClearCache(clientId string, callback func(confirmation *core.ClearCacheConfirmation, err error), props ...func(*core.ClearCacheRequest)) error {
	request := core.NewClearCacheRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*core.ClearCacheConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) DataTransfer(clientId string, callback func(confirmation *core.DataTransferConfirmation, err error), vendorId string, props ...func(request *core.DataTransferRequest)) error {
	request := core.NewDataTransferRequest(vendorId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*core.DataTransferConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) GetConfiguration(clientId string, callback func(confirmation *core.GetConfigurationConfirmation, err error), keys []string, props ...func(request *core.GetConfigurationRequest)) error {
	request := core.NewGetConfigurationRequest(keys)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*core.GetConfigurationConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) RemoteStartTransaction(clientId string, callback func(*core.RemoteStartTransactionConfirmation, error), idTag string, props ...func(*core.RemoteStartTransactionRequest)) error {
	request := core.NewRemoteStartTransactionRequest(idTag)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*core.RemoteStartTransactionConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) RemoteStopTransaction(clientId string, callback func(*core.RemoteStopTransactionConfirmation, error), transactionId int, props ...func(request *core.RemoteStopTransactionRequest)) error {
	request := core.NewRemoteStopTransactionRequest(transactionId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*core.RemoteStopTransactionConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) Reset(clientId string, callback func(*core.ResetConfirmation, error), resetType core.ResetType, props ...func(request *core.ResetRequest)) error {
	request := core.NewResetRequest(resetType)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*core.ResetConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) UnlockConnector(clientId string, callback func(*core.UnlockConnectorConfirmation, error), connectorId int, props ...func(*core.UnlockConnectorRequest)) error {
	request := core.NewUnlockConnectorRequest(connectorId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*core.UnlockConnectorConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) GetLocalListVersion(clientId string, callback func(*localauth.GetLocalListVersionConfirmation, error), props ...func(request *localauth.GetLocalListVersionRequest)) error {
	request := localauth.NewGetLocalListVersionRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*localauth.GetLocalListVersionConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) SendLocalList(clientId string, callback func(*localauth.SendLocalListConfirmation, error), version int, updateType localauth.UpdateType, props ...func(request *localauth.SendLocalListRequest)) error {
	request := localauth.NewSendLocalListRequest(version, updateType)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*localauth.SendLocalListConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) GetDiagnostics(clientId string, callback func(*firmware.GetDiagnosticsConfirmation, error), location string, props ...func(request *firmware.GetDiagnosticsRequest)) error {
	request := firmware.NewGetDiagnosticsRequest(location)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*firmware.GetDiagnosticsConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) UpdateFirmware(clientId string, callback func(*firmware.UpdateFirmwareConfirmation, error), location string, retrieveDate *types.DateTime, props ...func(request *firmware.UpdateFirmwareRequest)) error {
	request := firmware.NewUpdateFirmwareRequest(location, retrieveDate)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*firmware.UpdateFirmwareConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) ReserveNow(clientId string, callback func(*reservation.ReserveNowConfirmation, error), connectorId int, expiryDate *types.DateTime, idTag string, reservationId int, props ...func(request *reservation.ReserveNowRequest)) error {
	request := reservation.NewReserveNowRequest(connectorId, expiryDate, idTag, reservationId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*reservation.ReserveNowConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) CancelReservation(clientId string, callback func(*reservation.CancelReservationConfirmation, error), reservationId int, props ...func(request *reservation.CancelReservationRequest)) error {
	request := reservation.NewCancelReservationRequest(reservationId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*reservation.CancelReservationConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) TriggerMessage(clientId string, callback func(*remotetrigger.TriggerMessageConfirmation, error), requestedMessage remotetrigger.MessageTrigger, props ...func(request *remotetrigger.TriggerMessageRequest)) error {
	request := remotetrigger.NewTriggerMessageRequest(requestedMessage)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*remotetrigger.TriggerMessageConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) SetChargingProfile(clientId string, callback func(*smartcharging.SetChargingProfileConfirmation, error), connectorId int, chargingProfile *types.ChargingProfile, props ...func(request *smartcharging.SetChargingProfileRequest)) error {
	request := smartcharging.NewSetChargingProfileRequest(connectorId, chargingProfile)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*smartcharging.SetChargingProfileConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) ClearChargingProfile(clientId string, callback func(*smartcharging.ClearChargingProfileConfirmation, error), props ...func(request *smartcharging.ClearChargingProfileRequest)) error {
	request := smartcharging.NewClearChargingProfileRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*smartcharging.ClearChargingProfileConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) GetCompositeSchedule(clientId string, callback func(*smartcharging.GetCompositeScheduleConfirmation, error), connectorId int, duration int, props ...func(request *smartcharging.GetCompositeScheduleRequest)) error {
	request := smartcharging.NewGetCompositeScheduleRequest(connectorId, duration)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*smartcharging.GetCompositeScheduleConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) TriggerMessageExtended(clientId string, callback func(*extendedtriggermessage.ExtendedTriggerMessageResponse, error), requestedMessage extendedtriggermessage.ExtendedTriggerMessageType, props ...func(request *extendedtriggermessage.ExtendedTriggerMessageRequest)) error {
	request := extendedtriggermessage.NewExtendedTriggerMessageRequest(requestedMessage)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		confirmationCasted := confirmation.(*extendedtriggermessage.ExtendedTriggerMessageResponse)
		callback(confirmationCasted, protoError)
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)

}

func (cs *centralSystem) CertificateSigned(clientId string, callback func(*security.CertificateSignedResponse, error), csr string, props ...func(request *security.SignCertificateRequest)) error {
	request := security.NewSignCertificateRequest(csr)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		confirmationCasted := confirmation.(*security.CertificateSignedResponse)
		callback(confirmationCasted, protoError)
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) SignedUpdateFirmware(clientId string, callback func(*securefirmware.SignedUpdateFirmwareResponse, error), requestId int, firmware securefirmware.Firmware, props ...func(request *securefirmware.SignedUpdateFirmwareRequest)) error {
	request := securefirmware.NewSignedUpdateFirmwareRequest(requestId, firmware)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		confirmationCasted := confirmation.(*securefirmware.SignedUpdateFirmwareResponse)
		callback(confirmationCasted, protoError)
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) GetInstalledCertificateIds(clientId string, callback func(*certificates.GetInstalledCertificateIdsResponse, error), props ...func(request *certificates.GetInstalledCertificateIdsRequest)) error {
	request := certificates.NewGetInstalledCertificateIdsRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		confirmationCasted := confirmation.(*certificates.GetInstalledCertificateIdsResponse)
		callback(confirmationCasted, protoError)
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) InstallCertificate(clientId string, callback func(*certificates.InstallCertificateResponse, error), certificateType types.CertificateUse, certificate string, props ...func(request *certificates.InstallCertificateRequest)) error {
	request := certificates.NewInstallCertificateRequest(certificateType, certificate)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		confirmationCasted := confirmation.(*certificates.InstallCertificateResponse)
		callback(confirmationCasted, protoError)
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) DeleteCertificate(clientId string, callback func(*certificates.DeleteCertificateResponse, error), certificateHashData types.CertificateHashData, props ...func(request *certificates.DeleteCertificateRequest)) error {
	request := certificates.NewDeleteCertificateRequest(certificateHashData)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		confirmationCasted := confirmation.(*certificates.DeleteCertificateResponse)
		callback(confirmationCasted, protoError)
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) GetLog(clientId string, callback func(*logging.GetLogResponse, error), logType logging.LogType, requestID int, logParameters logging.LogParameters, props ...func(request *logging.GetLogRequest)) error {
	request := logging.NewGetLogRequest(logType, requestID, logParameters)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Response, protoError error) {
		confirmationCasted := confirmation.(*logging.GetLogResponse)
		callback(confirmationCasted, protoError)
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) SetSecurityHandler(handler security.CentralSystemHandler) {
	cs.securityHandler = handler
}

func (cs *centralSystem) SetLogHandler(handler logging.CentralSystemHandler) {
	cs.logHandler = handler
}

func (cs *centralSystem) SetSecureFirmwareHandler(handler securefirmware.CentralSystemHandler) {
	cs.secureFirmwareHandler = handler
}

func (cs *centralSystem) SetCoreHandler(handler core.CentralSystemHandler) {
	cs.coreHandler = handler
}

func (cs *centralSystem) SetLocalAuthListHandler(handler localauth.CentralSystemHandler) {
	cs.localAuthListHandler = handler
}

func (cs *centralSystem) SetFirmwareManagementHandler(handler firmware.CentralSystemHandler) {
	cs.firmwareHandler = handler
}

func (cs *centralSystem) SetReservationHandler(handler reservation.CentralSystemHandler) {
	cs.reservationHandler = handler
}

func (cs *centralSystem) SetRemoteTriggerHandler(handler remotetrigger.CentralSystemHandler) {
	cs.remoteTriggerHandler = handler
}

func (cs *centralSystem) SetSmartChargingHandler(handler smartcharging.CentralSystemHandler) {
	cs.smartChargingHandler = handler
}

func (cs *centralSystem) SetNewChargingStationValidationHandler(handler ws.CheckClientHandler) {
	cs.server.SetNewClientValidationHandler(handler)
}

func (cs *centralSystem) SetNewChargePointHandler(handler ChargePointConnectionHandler) {
	cs.server.SetNewClientHandler(func(chargePoint ws.Channel) {
		handler(chargePoint)
	})
}

func (cs *centralSystem) SetChargePointDisconnectedHandler(handler ChargePointConnectionHandler) {
	cs.server.SetDisconnectedClientHandler(func(chargePoint ws.Channel) {
		for cb, ok := cs.callbackQueue.Dequeue(chargePoint.ID()); ok; cb, ok = cs.callbackQueue.Dequeue(chargePoint.ID()) {
			err := ocpp.NewError(ocppj.GenericError, "client disconnected, no response received from client", "")
			cb(nil, err)
		}
		handler(chargePoint)
	})
}

func (cs *centralSystem) SendRequestAsync(clientId string, request ocpp.Request, callback func(confirmation ocpp.Response, err error)) error {
	featureName := request.GetFeatureName()
	if _, found := cs.server.GetProfileForFeature(featureName); !found {
		return fmt.Errorf("feature %v is unsupported on central system (missing profile), cannot send request", featureName)
	}
	switch featureName {
	case core.ChangeAvailabilityFeatureName, core.ChangeConfigurationFeatureName, core.ClearCacheFeatureName, core.DataTransferFeatureName, core.GetConfigurationFeatureName, core.RemoteStartTransactionFeatureName, core.RemoteStopTransactionFeatureName, core.ResetFeatureName, core.UnlockConnectorFeatureName,
		localauth.GetLocalListVersionFeatureName, localauth.SendLocalListFeatureName,
		firmware.GetDiagnosticsFeatureName, firmware.UpdateFirmwareFeatureName,
		reservation.ReserveNowFeatureName, reservation.CancelReservationFeatureName,
		remotetrigger.TriggerMessageFeatureName,
		smartcharging.SetChargingProfileFeatureName, smartcharging.ClearChargingProfileFeatureName, smartcharging.GetCompositeScheduleFeatureName,
		// OCPP 1.6 Extension supported messages
		security.CertificateSignedFeatureName,
		securefirmware.SignedUpdateFirmwareFeatureName,
		logging.GetLogFeatureName,
		extendedtriggermessage.ExtendedTriggerMessageFeatureName,
		certificates.GetInstalledCertificateIdsFeatureName, certificates.DeleteCertificateFeatureName, certificates.InstallCertificateFeatureName:
	default:
		return fmt.Errorf("unsupported action %v on central system, cannot send request", featureName)
	}

	send := func() error {
		return cs.server.SendRequest(clientId, request)
	}
	return cs.callbackQueue.TryQueue(clientId, send, callback)
}

func (cs *centralSystem) Start(listenPort int, listenPath string) {
	// Start server
	cs.server.Start(listenPort, listenPath)
}

func (cs *centralSystem) Stop() {
	cs.server.Stop()
}

func (cs *centralSystem) sendResponse(chargePointId string, confirmation ocpp.Response, err error, requestId string) {
	if err != nil {
		// Send error response
		if ocppError, ok := err.(*ocpp.Error); ok {
			err = cs.server.SendError(chargePointId, requestId, ocppError.Code, ocppError.Description, nil)
		} else {
			err = cs.server.SendError(chargePointId, requestId, ocppj.InternalError, err.Error(), nil)
		}
		if err != nil {
			// Error while sending an error. Will attempt to send a default error instead
			cs.server.HandleFailedResponseError(chargePointId, requestId, err, "")
			// Notify client implementation
			err = fmt.Errorf("error replying cp %s to request %s with 'internal error': %w", chargePointId, requestId, err)
			cs.error(err)
		}
		return
	}

	if confirmation == nil || reflect.ValueOf(confirmation).IsNil() {
		err = fmt.Errorf("empty confirmation to %s for request %s", chargePointId, requestId)
		// Sending a dummy error to server instead, then notify client implementation
		_ = cs.server.SendError(chargePointId, requestId, ocppj.GenericError, err.Error(), nil)
		cs.error(err)
		return
	}

	// send confirmation response
	err = cs.server.SendResponse(chargePointId, requestId, confirmation)
	if err != nil {
		// Error while sending an error. Will attempt to send a default error instead
		cs.server.HandleFailedResponseError(chargePointId, requestId, err, confirmation.GetFeatureName())
		// Notify client implementation
		err = fmt.Errorf("error replying cp %s to request %s: %w", chargePointId, requestId, err)
		cs.error(err)
	}
}

func (cs *centralSystem) notImplementedError(chargePointId string, requestId string, action string) {
	err := cs.server.SendError(chargePointId, requestId, ocppj.NotImplemented, fmt.Sprintf("no handler for action %v implemented", action), nil)
	if err != nil {
		err = fmt.Errorf("replying cp %s to request %s with 'not implemented': %w", chargePointId, requestId, err)
		cs.error(err)
	}
}

func (cs *centralSystem) notSupportedError(chargePointId string, requestId string, action string) {
	err := cs.server.SendError(chargePointId, requestId, ocppj.NotSupported, fmt.Sprintf("unsupported action %v on central system", action), nil)
	if err != nil {
		err = fmt.Errorf("replying cp %s to request %s with 'not supported': %w", chargePointId, requestId, err)
		cs.error(err)
	}
}

func (cs *centralSystem) handleIncomingRequest(chargePoint ChargePointConnection, request ocpp.Request, requestId string, action string) {
	profile, found := cs.server.GetProfileForFeature(action)
	// Check whether action is supported and a handler for it exists
	if !found {
		cs.notImplementedError(chargePoint.ID(), requestId, action)
		return
	} else {
		switch profile.Name {
		case core.ProfileName:
			if cs.coreHandler == nil {
				cs.notSupportedError(chargePoint.ID(), requestId, action)
				return
			}
		case localauth.ProfileName:
			if cs.localAuthListHandler == nil {
				cs.notSupportedError(chargePoint.ID(), requestId, action)
				return
			}
		case firmware.ProfileName:
			if cs.firmwareHandler == nil {
				cs.notSupportedError(chargePoint.ID(), requestId, action)
				return
			}
		case reservation.ProfileName:
			if cs.reservationHandler == nil {
				cs.notSupportedError(chargePoint.ID(), requestId, action)
				return
			}
		case remotetrigger.ProfileName:
			if cs.remoteTriggerHandler == nil {
				cs.notSupportedError(chargePoint.ID(), requestId, action)
				return
			}
		case smartcharging.ProfileName:
			if cs.smartChargingHandler == nil {
				cs.notSupportedError(chargePoint.ID(), requestId, action)
				return
			}
		case logging.ProfileName:
			if cs.logHandler == nil {
				cs.notSupportedError(chargePoint.ID(), requestId, action)
				return
			}
		case security.ProfileName:
			if cs.securityHandler == nil {
				cs.notSupportedError(chargePoint.ID(), requestId, action)
				return
			}
		case securefirmware.ProfileName:
			if cs.secureFirmwareHandler == nil {
				cs.notSupportedError(chargePoint.ID(), requestId, action)
				return
			}
		}
	}
	var confirmation ocpp.Response
	var err error
	// Execute in separate goroutine, so the caller goroutine is available
	go func() {
		switch action {
		case core.BootNotificationFeatureName:
			confirmation, err = cs.coreHandler.OnBootNotification(chargePoint.ID(), request.(*core.BootNotificationRequest))
		case core.AuthorizeFeatureName:
			confirmation, err = cs.coreHandler.OnAuthorize(chargePoint.ID(), request.(*core.AuthorizeRequest))
		case core.DataTransferFeatureName:
			confirmation, err = cs.coreHandler.OnDataTransfer(chargePoint.ID(), request.(*core.DataTransferRequest))
		case core.HeartbeatFeatureName:
			confirmation, err = cs.coreHandler.OnHeartbeat(chargePoint.ID(), request.(*core.HeartbeatRequest))
		case core.MeterValuesFeatureName:
			confirmation, err = cs.coreHandler.OnMeterValues(chargePoint.ID(), request.(*core.MeterValuesRequest))
		case core.StartTransactionFeatureName:
			confirmation, err = cs.coreHandler.OnStartTransaction(chargePoint.ID(), request.(*core.StartTransactionRequest))
		case core.StopTransactionFeatureName:
			confirmation, err = cs.coreHandler.OnStopTransaction(chargePoint.ID(), request.(*core.StopTransactionRequest))
		case core.StatusNotificationFeatureName:
			confirmation, err = cs.coreHandler.OnStatusNotification(chargePoint.ID(), request.(*core.StatusNotificationRequest))
		case firmware.DiagnosticsStatusNotificationFeatureName:
			confirmation, err = cs.firmwareHandler.OnDiagnosticsStatusNotification(chargePoint.ID(), request.(*firmware.DiagnosticsStatusNotificationRequest))
		case firmware.FirmwareStatusNotificationFeatureName:
			confirmation, err = cs.firmwareHandler.OnFirmwareStatusNotification(chargePoint.ID(), request.(*firmware.FirmwareStatusNotificationRequest))
		case security.SignCertificateFeatureName:
			confirmation, err = cs.securityHandler.OnSignCertificate(chargePoint.ID(), request.(*security.SignCertificateRequest))
		case security.SecurityEventNotificationFeatureName:
			confirmation, err = cs.securityHandler.OnSecurityEventNotification(chargePoint.ID(), request.(*security.SecurityEventNotificationRequest))
		case logging.LogStatusNotificationFeatureName:
			confirmation, err = cs.logHandler.OnLogStatusNotification(chargePoint.ID(), request.(*logging.LogStatusNotificationRequest))
		case securefirmware.SignedFirmwareStatusNotificationFeatureName:
			confirmation, err = cs.secureFirmwareHandler.OnSignedFirmwareStatusNotification(chargePoint.ID(), request.(*securefirmware.SignedFirmwareStatusNotificationRequest))
		default:
			cs.notSupportedError(chargePoint.ID(), requestId, action)
			return
		}
		cs.sendResponse(chargePoint.ID(), confirmation, err, requestId)
	}()
}

func (cs *centralSystem) handleIncomingConfirmation(chargePoint ChargePointConnection, confirmation ocpp.Response, requestId string) {
	if callback, ok := cs.callbackQueue.Dequeue(chargePoint.ID()); ok {
		// Execute in separate goroutine, so the caller goroutine is available
		go callback(confirmation, nil)
	} else {
		err := fmt.Errorf("no handler available for call of type %v from client %s for request %s", confirmation.GetFeatureName(), chargePoint.ID(), requestId)
		cs.error(err)
	}
}

func (cs *centralSystem) handleIncomingError(chargePoint ChargePointConnection, err *ocpp.Error, details interface{}) {
	if callback, ok := cs.callbackQueue.Dequeue(chargePoint.ID()); ok {
		// Execute in separate goroutine, so the caller goroutine is available
		go callback(nil, err)
	} else {
		err := fmt.Errorf("no handler available for call error %w from client %s", err, chargePoint.ID())
		cs.error(err)
	}
}

func (cs *centralSystem) handleCanceledRequest(chargePointID string, request ocpp.Request, err *ocpp.Error) {
	if callback, ok := cs.callbackQueue.Dequeue(chargePointID); ok {
		// Execute in separate goroutine, so the caller goroutine is available
		go callback(nil, err)
	} else {
		err := fmt.Errorf("no handler available for canceled request %s for client %s: %w",
			request.GetFeatureName(), chargePointID, err)
		cs.error(err)
	}
}
