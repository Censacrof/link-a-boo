import { Construct } from "constructs";
import * as s3 from "aws-cdk-lib/aws-s3";
import * as s3deploy from "aws-cdk-lib/aws-s3-deployment";
import * as cdk from "aws-cdk-lib/core";

export class LinkABooRedoc extends Construct {
  public readonly bucket: s3.Bucket;

  constructor(scope: Construct, id: string) {
    super(scope, id);

    this.bucket = new s3.Bucket(this, "RedocBucket", {
      publicReadAccess: true,
      blockPublicAccess: new s3.BlockPublicAccess({
        blockPublicAcls: false,
        ignorePublicAcls: false,
        blockPublicPolicy: false,
        restrictPublicBuckets: false,
      }),
      websiteIndexDocument: "index.html",
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });

    new s3deploy.BucketDeployment(this, "DeployRedoc", {
      sources: [s3deploy.Source.asset("./redoc")],
      destinationBucket: this.bucket,
    });

    new cdk.CfnOutput(this, "RedocURL", {
      value: this.bucket.bucketWebsiteUrl,
    });
  }
}
