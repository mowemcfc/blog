import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda'
import * as apigw from 'aws-cdk-lib/aws-apigateway'


export class DeployStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
      const htmxGoLambda = new lambda.Function(this, 'HtmxGoLambda', {
        runtime: lambda.Runtime.PROVIDED_AL2023,
        handler: 'main',
        code: lambda.Code.fromAsset('../blog.zip'),
        environment: {
          // Add any environment variables your Lambda function needs
        }
      });

      // Define API Gateway endpoint
      const api = new apigw.RestApi(this, 'HtmxGoLambdaAPI');
      const htmxGoLambdaIntegration = new apigw.LambdaIntegration(htmxGoLambda);
      api.root.addMethod('GET', htmxGoLambdaIntegration);
  }
}

