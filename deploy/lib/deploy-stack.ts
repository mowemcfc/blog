import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as apigw from "aws-cdk-lib/aws-apigateway";
import * as acm from "aws-cdk-lib/aws-certificatemanager";
import * as route53 from "aws-cdk-lib/aws-route53";
import * as targets from "aws-cdk-lib/aws-route53-targets";

export class DeployStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
    const htmxGoLambda = new lambda.Function(this, "HtmxGoLambda", {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: "main",
      code: lambda.Code.fromAsset("../blog.zip"),
      environment: {},
    });

    const jcartershHostedZone = route53.HostedZone.fromHostedZoneAttributes(
      this,
      "jcartershHostedZone",
      {
        hostedZoneId: "Z1033121V15B6T3SS84I",
        zoneName: "jcarter.sh",
      },
    );

    const certificate = new acm.Certificate(this, "HtmxGoAPICertificate", {
      domainName: "jcarter.sh",
      subjectAlternativeNames: ["*.jcarter.sh"],
      validation: acm.CertificateValidation.fromDns(jcartershHostedZone),
    });

    const api = new apigw.LambdaRestApi(this, "HtmxGoLambdaAPI", {
      handler: htmxGoLambda,
      deploy: true,
      proxy: true,
      domainName: {
        domainName: "jcarter.sh",
        certificate: certificate,
      },
    });

    const record = new route53.ARecord(this, "ARecord", {
      target: route53.RecordTarget.fromAlias(new targets.ApiGateway(api)),
      zone: jcartershHostedZone,
    });
  }
}
