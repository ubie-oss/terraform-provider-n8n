package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func NewCredentialSMTPResource() resource.Resource {
	return newTypedCredentialResource(smtpSpec())
}

func NewCredentialAWSResource() resource.Resource {
	return newTypedCredentialResource(awsSpec())
}

func smtpSpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_smtp",
		N8nType:         "smtp",
		DocFile:         "internal/provider/docs/resources/credential_smtp.md",
		ExtraAttributes: map[string]schema.Attribute{
			"user":              stringState("SMTP username / mailbox address.", true),
			"password":          stringWriteOnly("SMTP password or app password. Write-only; never stored in state.", true),
			"host":              stringState("SMTP host, for example `smtp.example.com`.", true),
			"port":              int64State("SMTP port. n8n defaults to 465 when omitted.", false),
			"secure":            boolState("When true, use implicit SSL/TLS (typical for port 465)."),
			"disable_start_tls": boolState("When SSL/TLS is off, prevent opportunistic STARTTLS upgrades."),
			"host_name":         stringState("Client host name (FQDN) sent to the SMTP server. Leave empty unless your provider requires it.", false),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "user", bag.strings["user"]); err != nil {
				return nil, err
			}
			if err := bagPutRequiredString(m, "password", bag.strings["password"]); err != nil {
				return nil, err
			}
			if err := bagPutRequiredString(m, "host", bag.strings["host"]); err != nil {
				return nil, err
			}
			bagPutInt64(m, "port", bag.int64s["port"])
			bagPutBool(m, "secure", bag.bools["secure"])
			bagPutBool(m, "disableStartTls", bag.bools["disable_start_tls"])
			bagPutString(m, "hostName", bag.strings["host_name"])
			return m, nil
		},
	}
}

func awsSpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_aws",
		N8nType:         "aws",
		DocFile:         "internal/provider/docs/resources/credential_aws.md",
		ExtraAttributes: map[string]schema.Attribute{
			"region":                   stringState("AWS region, for example `us-east-1`.", true),
			"access_key_id":            stringState("IAM access key id.", true),
			"secret_access_key":        stringWriteOnly("IAM secret access key. Write-only; never stored in state.", true),
			"temporary_credentials":    boolState("When true, send session_token for STS temporary credentials."),
			"session_token":            stringWriteOnly("STS session token used when temporary_credentials is true. Write-only; never stored in state.", false),
			"custom_endpoints":         boolState("When true, send VPC custom service endpoints."),
			"rekognition_endpoint":     stringState("Custom Rekognition endpoint.", false),
			"lambda_endpoint":          stringState("Custom Lambda endpoint.", false),
			"sns_endpoint":             stringState("Custom SNS endpoint.", false),
			"ses_endpoint":             stringState("Custom SES endpoint.", false),
			"sqs_endpoint":             stringState("Custom SQS endpoint.", false),
			"s3_endpoint":              stringState("Custom S3 endpoint.", false),
			"ssm_endpoint":             stringState("Custom SSM endpoint.", false),
			"bedrock_endpoint":         stringState("Custom Bedrock control-plane endpoint.", false),
			"bedrock_runtime_endpoint": stringState("Custom Bedrock runtime endpoint.", false),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "region", bag.strings["region"]); err != nil {
				return nil, err
			}
			if err := bagPutRequiredString(m, "accessKeyId", bag.strings["access_key_id"]); err != nil {
				return nil, err
			}
			if err := bagPutRequiredString(m, "secretAccessKey", bag.strings["secret_access_key"]); err != nil {
				return nil, err
			}
			bagPutBool(m, "temporaryCredentials", bag.bools["temporary_credentials"])
			bagPutString(m, "sessionToken", bag.strings["session_token"])
			bagPutBool(m, "customEndpoints", bag.bools["custom_endpoints"])
			bagPutString(m, "rekognitionEndpoint", bag.strings["rekognition_endpoint"])
			bagPutString(m, "lambdaEndpoint", bag.strings["lambda_endpoint"])
			bagPutString(m, "snsEndpoint", bag.strings["sns_endpoint"])
			bagPutString(m, "sesEndpoint", bag.strings["ses_endpoint"])
			bagPutString(m, "sqsEndpoint", bag.strings["sqs_endpoint"])
			bagPutString(m, "s3Endpoint", bag.strings["s3_endpoint"])
			bagPutString(m, "ssmEndpoint", bag.strings["ssm_endpoint"])
			bagPutString(m, "bedrockEndpoint", bag.strings["bedrock_endpoint"])
			bagPutString(m, "bedrockRuntimeEndpoint", bag.strings["bedrock_runtime_endpoint"])
			return m, nil
		},
	}
}
