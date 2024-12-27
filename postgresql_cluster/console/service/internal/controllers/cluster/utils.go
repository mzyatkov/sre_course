package cluster

import (
	"context"
	"encoding/json"
	"postgresql-cluster-console/internal/storage"
	"postgresql-cluster-console/models"
	"postgresql-cluster-console/pkg/tracer"

	"github.com/rs/zerolog"
)

type ParamLocation uint8

const (
	UnknownParamLocation   ParamLocation = 0
	EnvParamLocation       ParamLocation = 1
	ExtraVarsParamLocation ParamLocation = 2
)

func getSecretEnvs(ctx context.Context, log zerolog.Logger, db storage.IStorage, secretID int64, secretKey string) ([]string, ParamLocation, error) {
	localLog := log.With().Str("cid", ctx.Value(tracer.CtxCidKey{}).(string)).Logger()
	secretView, err := db.GetSecret(ctx, secretID)
	if err != nil {
		return nil, UnknownParamLocation, err
	}
	localLog.Trace().Any("secret_view", secretView).Msg("got secret view from db")
	secretVal, err := db.GetSecretVal(ctx, secretID, secretKey)
	if err != nil {
		return nil, UnknownParamLocation, err
	}
	localLog.Trace().Msgf("secretVal %s", string(secretVal))

	switch models.SecretType(secretView.Type) {
	case models.SecretTypeAws:
		var sec models.RequestSecretValueAws
		err = json.Unmarshal(secretVal, &sec)
		if err != nil {
			return nil, UnknownParamLocation, err
		}

		return []string{"AWS_ACCESS_KEY_ID=" + sec.AWSACCESSKEYID, "AWS_SECRET_ACCESS_KEY=" + sec.AWSSECRETACCESSKEY}, EnvParamLocation, nil
	case models.SecretTypeGcp:
		var sec models.RequestSecretValueGcp
		err = json.Unmarshal(secretVal, &sec)
		if err != nil {
			return nil, UnknownParamLocation, err
		}

		return []string{"GCP_SERVICE_ACCOUNT_CONTENTS=" + sec.GCPSERVICEACCOUNTCONTENTS}, EnvParamLocation, nil
	case models.SecretTypeAzure:
		var sec models.RequestSecretValueAzure
		err = json.Unmarshal(secretVal, &sec)
		if err != nil {
			return nil, UnknownParamLocation, err
		}

		return []string{
			"AZURE_SUBSCRIPTION_ID=" + sec.AZURESUBSCRIPTIONID,
			"AZURE_CLIENT_ID=" + sec.AZURECLIENTID,
			"AZURE_SECRET=" + sec.AZURESECRET,
			"AZURE_TENANT=" + sec.AZURETENANT,
		}, EnvParamLocation, nil
	case models.SecretTypeDigitalocean:
		var sec models.RequestSecretValueDigitalOcean
		err = json.Unmarshal(secretVal, &sec)
		if err != nil {
			return nil, UnknownParamLocation, err
		}

		return []string{"DO_API_TOKEN=" + sec.DOAPITOKEN}, EnvParamLocation, nil
	case models.SecretTypeHetzner:
		var sec models.RequestSecretValueHetzner
		err = json.Unmarshal(secretVal, &sec)
		if err != nil {
			return nil, UnknownParamLocation, err
		}

		return []string{"HCLOUD_API_TOKEN=" + sec.HCLOUDAPITOKEN}, EnvParamLocation, nil
	case models.SecretTypeSSHKey:
		var sec models.RequestSecretValueSSHKey
		err = json.Unmarshal(secretVal, &sec)
		if err != nil {
			return nil, UnknownParamLocation, err
		}

		return []string{"SSH_PRIVATE_KEY_CONTENT=" + sec.SSHPRIVATEKEY}, EnvParamLocation, nil
	case models.SecretTypePassword:
		var sec models.RequestSecretValuePassword
		err = json.Unmarshal(secretVal, &sec)
		if err != nil {
			return nil, UnknownParamLocation, err
		}

		return []string{"ansible_user=" + sec.USERNAME, "ansible_ssh_pass=" + sec.PASSWORD, "ansible_sudo_pass=" + sec.PASSWORD}, ExtraVarsParamLocation, nil
	default:
		return nil, UnknownParamLocation, nil
	}
}
