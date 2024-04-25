import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as apigw from "aws-cdk-lib/aws-apigateway";
import * as acm from "aws-cdk-lib/aws-certificatemanager";
import * as route53 from "aws-cdk-lib/aws-route53";
import * as targets from "aws-cdk-lib/aws-route53-targets";
import * as s3 from "aws-cdk-lib/aws-s3";
import * as s3deploy from "aws-cdk-lib/aws-s3-deployment";
import * as cloudfront from "aws-cdk-lib/aws-cloudfront";
import * as origins from "aws-cdk-lib/aws-cloudfront-origins";

export class DeployStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
    const htmxGoLambda = new lambda.Function(this, "HtmxGoLambda", {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: "main",
      code: lambda.Code.fromAsset("../blog.zip"),
      environment: {
        // Add any environment variables your Lambda function needs
      },
    });
    const assetsBucket = new s3.Bucket(this, "BlogAssetBucket", {});
    new s3deploy.BucketDeployment(this, "BlogDeployFiles", {
      sources: [s3deploy.Source.asset("../static")],
      destinationBucket: assetsBucket,
    });

    const jcartershHostedZone = route53.HostedZone.fromHostedZoneAttributes(
      this,
      "jcartershHostedZone",
      {
        hostedZoneId: "Z1033121V15B6T3SS84I",
        zoneName: "jcarter.sh",
      },
    );

    const globalCertificate = acm.Certificate.fromCertificateArn(
      this,
      "Certificate",
      "arn:aws:acm:us-east-1:891854796411:certificate/e9a5f789-4e11-4601-ae16-425b6a6e72c4",
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

    const originAccessIdentity = new cloudfront.OriginAccessIdentity(
      this,
      "OriginAccessIdentity",
    );
    assetsBucket.grantRead(originAccessIdentity);
    const bucketOrigin = new origins.S3Origin(assetsBucket, {
      originAccessIdentity,
    });
    const apigwOrigin = new origins.RestApiOrigin(api);
    const distribution = new cloudfront.Distribution(this, "BlogDistribution", {
      defaultBehavior: {
        origin: apigwOrigin,
      },
      additionalBehaviors: {
        "/js/*": { origin: bucketOrigin },
        "/css/*": { origin: bucketOrigin },
        "/images/*": { origin: bucketOrigin },
      },
      certificate: globalCertificate,
      domainNames: ["jcarter.sh"],
    });

    const record = new route53.ARecord(this, "ARecord", {
      target: route53.RecordTarget.fromAlias(
        new targets.CloudFrontTarget(distribution),
      ),
      zone: jcartershHostedZone,
    });
  }
}
