// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appsync_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appsync"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfappsync "github.com/hashicorp/terraform-provider-aws/internal/service/appsync"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func testAccResolver_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var resolver1 appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, appsync.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.AppSyncServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResolverDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResolverConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(ctx, resourceName, &resolver1),
					acctest.MatchResourceAttrRegionalARN(resourceName, names.AttrARN, "appsync", regexache.MustCompile("apis/.+/types/.+/resolvers/.+")),
					resource.TestCheckResourceAttr(resourceName, "data_source", rName),
					resource.TestCheckResourceAttrSet(resourceName, "request_template"),
					resource.TestCheckResourceAttr(resourceName, "max_batch_size", acctest.CtZero),
					resource.TestCheckResourceAttr(resourceName, "sync_config.#", acctest.CtZero),
					resource.TestCheckResourceAttr(resourceName, "runtime.#", acctest.CtZero),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResolver_code(t *testing.T) {
	ctx := acctest.Context(t)
	var resolver1 appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, appsync.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.AppSyncServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResolverDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResolverConfig_code(rName, "test-fixtures/test-code.js"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(ctx, resourceName, &resolver1),
					resource.TestCheckResourceAttr(resourceName, "runtime.#", acctest.CtOne),
					resource.TestCheckResourceAttr(resourceName, "runtime.0.name", "APPSYNC_JS"),
					resource.TestCheckResourceAttr(resourceName, "runtime.0.runtime_version", "1.0.0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResolver_syncConfig(t *testing.T) {
	ctx := acctest.Context(t)
	var resolver1 appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, appsync.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.AppSyncServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResolverDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResolverConfig_sync(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(ctx, resourceName, &resolver1),
					resource.TestCheckResourceAttr(resourceName, "sync_config.#", acctest.CtOne),
					resource.TestCheckResourceAttr(resourceName, "sync_config.0.conflict_detection", "VERSION"),
					resource.TestCheckResourceAttr(resourceName, "sync_config.0.conflict_handler", "OPTIMISTIC_CONCURRENCY"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResolver_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var api1 appsync.GraphqlApi
	var resolver1 appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	appsyncGraphqlApiResourceName := "aws_appsync_graphql_api.test"
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, appsync.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.AppSyncServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResolverDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResolverConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGraphQLAPIExists(ctx, appsyncGraphqlApiResourceName, &api1),
					testAccCheckResolverExists(ctx, resourceName, &resolver1),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfappsync.ResourceResolver(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccResolver_dataSource(t *testing.T) {
	ctx := acctest.Context(t)
	var resolver1, resolver2 appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, appsync.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.AppSyncServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResolverDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResolverConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(ctx, resourceName, &resolver1),
					resource.TestCheckResourceAttr(resourceName, "data_source", rName),
				),
			},
			{
				Config: testAccResolverConfig_dataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(ctx, resourceName, &resolver2),
					resource.TestCheckResourceAttr(resourceName, "data_source", "test_ds_2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResolver_DataSource_lambda(t *testing.T) {
	ctx := acctest.Context(t)
	var resolver appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, appsync.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.AppSyncServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResolverDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResolverConfig_dataSourceLambda(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(ctx, resourceName, &resolver),
					resource.TestCheckResourceAttr(resourceName, "data_source", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResolver_requestTemplate(t *testing.T) {
	ctx := acctest.Context(t)
	var resolver1, resolver2 appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, appsync.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.AppSyncServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResolverDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResolverConfig_requestTemplate(rName, "/"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(ctx, resourceName, &resolver1),
					resource.TestMatchResourceAttr(resourceName, "request_template", regexache.MustCompile("resourcePath\": \"/\"")),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResolverConfig_requestTemplate(rName, "/test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(ctx, resourceName, &resolver2),
					resource.TestMatchResourceAttr(resourceName, "request_template", regexache.MustCompile("resourcePath\": \"/test\"")),
				),
			},
		},
	})
}

func testAccResolver_responseTemplate(t *testing.T) {
	ctx := acctest.Context(t)
	var resolver1, resolver2 appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, appsync.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.AppSyncServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResolverDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResolverConfig_responseTemplate(rName, 200),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(ctx, resourceName, &resolver1),
					resource.TestMatchResourceAttr(resourceName, "response_template", regexache.MustCompile(`ctx\.result\.statusCode == 200`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResolverConfig_responseTemplate(rName, 201),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(ctx, resourceName, &resolver2),
					resource.TestMatchResourceAttr(resourceName, "response_template", regexache.MustCompile(`ctx\.result\.statusCode == 201`)),
				),
			},
		},
	})
}

func testAccResolver_multipleResolvers(t *testing.T) {
	ctx := acctest.Context(t)
	var resolver appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, appsync.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.AppSyncServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResolverDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResolverConfig_multiple(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(ctx, resourceName+acctest.CtOne, &resolver),
					testAccCheckResolverExists(ctx, resourceName+acctest.CtTwo, &resolver),
					testAccCheckResolverExists(ctx, resourceName+"3", &resolver),
					testAccCheckResolverExists(ctx, resourceName+"4", &resolver),
					testAccCheckResolverExists(ctx, resourceName+"5", &resolver),
					testAccCheckResolverExists(ctx, resourceName+"6", &resolver),
					testAccCheckResolverExists(ctx, resourceName+"7", &resolver),
					testAccCheckResolverExists(ctx, resourceName+"8", &resolver),
					testAccCheckResolverExists(ctx, resourceName+"9", &resolver),
					testAccCheckResolverExists(ctx, resourceName+"10", &resolver),
				),
			},
		},
	})
}

func testAccResolver_pipeline(t *testing.T) {
	ctx := acctest.Context(t)
	var resolver appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, appsync.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.AppSyncServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResolverDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResolverConfig_pipeline(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(ctx, resourceName, &resolver),
					resource.TestCheckResourceAttr(resourceName, "pipeline_config.0.functions.#", acctest.CtOne),
					resource.TestCheckResourceAttrPair(resourceName, "pipeline_config.0.functions.0", "aws_appsync_function.test", "function_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResolver_caching(t *testing.T) {
	ctx := acctest.Context(t)
	var resolver appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, appsync.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.AppSyncServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResolverDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResolverConfig_caching(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(ctx, resourceName, &resolver),
					resource.TestCheckResourceAttr(resourceName, "caching_config.0.caching_keys.#", acctest.CtTwo),
					resource.TestCheckResourceAttr(resourceName, "caching_config.0.ttl", "60"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckResolverDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).AppSyncConn(ctx)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_appsync_resolver" {
				continue
			}

			apiID, typeName, fieldName, err := tfappsync.DecodeResolverID(rs.Primary.ID)

			if err != nil {
				return err
			}

			input := &appsync.GetResolverInput{
				ApiId:     aws.String(apiID),
				TypeName:  aws.String(typeName),
				FieldName: aws.String(fieldName),
			}

			_, err = conn.GetResolverWithContext(ctx, input)

			if tfawserr.ErrCodeEquals(err, appsync.ErrCodeNotFoundException) {
				continue
			}

			if err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckResolverExists(ctx context.Context, name string, resolver *appsync.Resolver) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource has no ID: %s", name)
		}

		apiID, typeName, fieldName, err := tfappsync.DecodeResolverID(rs.Primary.ID)

		if err != nil {
			return err
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).AppSyncConn(ctx)

		input := &appsync.GetResolverInput{
			ApiId:     aws.String(apiID),
			TypeName:  aws.String(typeName),
			FieldName: aws.String(fieldName),
		}

		output, err := conn.GetResolverWithContext(ctx, input)

		if err != nil {
			return err
		}

		*resolver = *output.Resolver

		return nil
	}
}

func testAccResolverConfig_base(rName string) string {
	return fmt.Sprintf(`
resource "aws_appsync_graphql_api" "test" {
  authentication_type = "API_KEY"
  name                = %[1]q

  schema = <<EOF
type Mutation {
	putPost(id: ID!, title: String!): Post
}

type Post {
	id: ID!
	title: String!
}

type Query {
	singlePost(id: ID!): Post
}

schema {
	query: Query
	mutation: Mutation
}
EOF
}

resource "aws_appsync_datasource" "test" {
  api_id = aws_appsync_graphql_api.test.id
  name   = %[1]q
  type   = "HTTP"

  http_config {
    endpoint = "http://example.com"
  }
}
`, rName)
}

func testAccResolverConfig_basic(rName string) string {
	return testAccResolverConfig_base(rName) + `
resource "aws_appsync_resolver" "test" {
  api_id      = aws_appsync_graphql_api.test.id
  field       = "singlePost"
  type        = "Query"
  data_source = aws_appsync_datasource.test.name

  request_template = <<EOF
{
    "version": "2018-05-29",
    "method": "GET",
    "resourcePath": "/",
    "params":{
        "headers": $utils.http.copyheaders($ctx.request.headers)
    }
}
EOF

  response_template = <<EOF
#if($ctx.result.statusCode == 200)
    $ctx.result.body
#else
    $utils.appendError($ctx.result.body, $ctx.result.statusCode)
#end
EOF
}
`
}

func testAccResolverConfig_dataSource(rName string) string {
	return testAccResolverConfig_base(rName) + `
resource "aws_appsync_datasource" "test2" {
  api_id = aws_appsync_graphql_api.test.id
  name   = "test_ds_2"
  type   = "HTTP"

  http_config {
    endpoint = "http://example.com"
  }
}

resource "aws_appsync_resolver" "test" {
  api_id      = aws_appsync_graphql_api.test.id
  field       = "singlePost"
  type        = "Query"
  data_source = aws_appsync_datasource.test2.name

  request_template = <<EOF
{
    "version": "2018-05-29",
    "method": "GET",
    "resourcePath": "/",
    "params":{
        "headers": $utils.http.copyheaders($ctx.request.headers)
    }
}
EOF

  response_template = <<EOF
#if($ctx.result.statusCode == 200)
    $ctx.result.body
#else
    $utils.appendError($ctx.result.body, $ctx.result.statusCode)
#end
EOF
}
`
}

func testAccResolverConfig_dataSourceLambda(rName string) string {
	return testAccDatasourceConfig_baseLambda(rName) + fmt.Sprintf(`
resource "aws_appsync_graphql_api" "test" {
  authentication_type = "API_KEY"
  name                = %q

  schema = <<EOF
type Mutation {
	putPost(id: ID!, title: String!): Post
}

type Post {
	id: ID!
	title: String!
}

type Query {
	singlePost(id: ID!): Post
}

schema {
	query: Query
	mutation: Mutation
}
EOF
}

resource "aws_appsync_datasource" "test" {
  api_id           = aws_appsync_graphql_api.test.id
  name             = %q
  service_role_arn = aws_iam_role.test.arn
  type             = "AWS_LAMBDA"

  lambda_config {
    function_arn = aws_lambda_function.test.arn
  }
}

resource "aws_appsync_resolver" "test" {
  api_id      = aws_appsync_graphql_api.test.id
  field       = "singlePost"
  type        = "Query"
  data_source = aws_appsync_datasource.test.name
}
`, rName, rName)
}

func testAccResolverConfig_requestTemplate(rName, resourcePath string) string {
	return testAccResolverConfig_base(rName) + fmt.Sprintf(`
resource "aws_appsync_resolver" "test" {
  api_id      = aws_appsync_graphql_api.test.id
  field       = "singlePost"
  type        = "Query"
  data_source = aws_appsync_datasource.test.name

  request_template = <<EOF
{
    "version": "2018-05-29",
    "method": "GET",
    "resourcePath": %[1]q,
    "params":{
        "headers": $utils.http.copyheaders($ctx.request.headers)
    }
}
EOF

  response_template = <<EOF
#if($ctx.result.statusCode == 200)
    $ctx.result.body
#else
    $utils.appendError($ctx.result.body, $ctx.result.statusCode)
#end
EOF
}
`, resourcePath)
}

func testAccResolverConfig_responseTemplate(rName string, statusCode int) string {
	return testAccResolverConfig_base(rName) + fmt.Sprintf(`
resource "aws_appsync_resolver" "test" {
  api_id      = aws_appsync_graphql_api.test.id
  field       = "singlePost"
  type        = "Query"
  data_source = aws_appsync_datasource.test.name

  request_template = <<EOF
{
    "version": "2018-05-29",
    "method": "GET",
    "resourcePath": "/",
    "params":{
        ## you can forward the headers using the below utility
        "headers": $utils.http.copyheaders($ctx.request.headers)
    }
}
EOF

  response_template = <<EOF
#if($ctx.result.statusCode == %[1]d)
    $ctx.result.body
#else
    $utils.appendError($ctx.result.body, $ctx.result.statusCode)
#end
EOF
}
`, statusCode)
}

func testAccResolverConfig_multiple(rName string) string {
	var queryFields string
	var resolverResources string
	for i := 1; i <= 10; i++ {
		queryFields = queryFields + fmt.Sprintf(`
	singlePost%d(id: ID!): Post
`, i)
		resolverResources = resolverResources + fmt.Sprintf(`
resource "aws_appsync_resolver" "test%d" {
  api_id      = aws_appsync_graphql_api.test.id
  field       = "singlePost%d"
  type        = "Query"
  data_source = aws_appsync_datasource.test.name

  request_template = <<EOF
{
    "version": "2018-05-29",
    "method": "GET",
    "resourcePath": "/",
    "params":{
        "headers": $utils.http.copyheaders($ctx.request.headers)
    }
}
EOF

  response_template = <<EOF
#if($ctx.result.statusCode == 200)
    $ctx.result.body
#else
    $utils.appendError($ctx.result.body, $ctx.result.statusCode)
#end
EOF
}
`, i, i)
	}

	return fmt.Sprintf(`
resource "aws_appsync_graphql_api" "test" {
  authentication_type = "API_KEY"
  name                = %q

  schema = <<EOF
type Mutation {
	putPost(id: ID!, title: String!): Post
}

type Post {
	id: ID!
	title: String!
}

type Query {
%s
}

schema {
	query: Query
	mutation: Mutation
}
EOF
}

resource "aws_appsync_datasource" "test" {
  api_id = aws_appsync_graphql_api.test.id
  name   = %q
  type   = "HTTP"

  http_config {
    endpoint = "http://example.com"
  }
}

%s

`, rName, queryFields, rName, resolverResources)
}

func testAccResolverConfig_pipeline(rName string) string {
	return testAccResolverConfig_base(rName) + fmt.Sprintf(`
resource "aws_appsync_function" "test" {
  api_id                   = aws_appsync_graphql_api.test.id
  data_source              = aws_appsync_datasource.test.name
  name                     = %[1]q
  request_mapping_template = <<EOF
{
		"version": "2018-05-29",
		"method": "GET",
		"resourcePath": "/",
		"params":{
				"headers": $utils.http.copyheaders($ctx.request.headers)
		}
}
EOF

  response_mapping_template = <<EOF
#if($ctx.result.statusCode == 200)
		$ctx.result.body
#else
		$utils.appendError($ctx.result.body, $ctx.result.statusCode)
#end
EOF
}

resource "aws_appsync_resolver" "test" {
  api_id           = aws_appsync_graphql_api.test.id
  field            = "singlePost"
  type             = "Query"
  kind             = "PIPELINE"
  request_template = <<EOF
{
		"version": "2018-05-29",
		"method": "GET",
		"resourcePath": "/",
		"params":{
				"headers": $utils.http.copyheaders($ctx.request.headers)
		}
}
EOF

  response_template = <<EOF
#if($ctx.result.statusCode == 200)
		$ctx.result.body
#else
		$utils.appendError($ctx.result.body, $ctx.result.statusCode)
#end
EOF

  pipeline_config {
    functions = [aws_appsync_function.test.function_id]
  }
}

`, rName)
}

func testAccResolverConfig_caching(rName string) string {
	return testAccResolverConfig_base(rName) + `
resource "aws_appsync_resolver" "test" {
  api_id           = aws_appsync_graphql_api.test.id
  field            = "singlePost"
  type             = "Query"
  kind             = "UNIT"
  data_source      = aws_appsync_datasource.test.name
  request_template = <<EOF
{
    "version": "2018-05-29",
    "method": "GET",
    "resourcePath": "/",
    "params":{
        "headers": $utils.http.copyheaders($ctx.request.headers)
    }
}
EOF

  response_template = <<EOF
#if($ctx.result.statusCode == 200)
    $ctx.result.body
#else
    $utils.appendError($ctx.result.body, $ctx.result.statusCode)
#end
EOF

  caching_config {
    caching_keys = [
      "$context.identity.sub",
      "$context.arguments.id",
    ]
    ttl = 60
  }
}
`
}

func testAccResolverConfig_sync(rName string) string {
	return testAccDatasourceConfig_baseDynamoDB(rName) + fmt.Sprintf(`
resource "aws_appsync_graphql_api" "test" {
  authentication_type = "API_KEY"
  name                = %[1]q

  schema = <<EOF
type Mutation {
	putPost(id: ID!, title: String!): Post
}

type Post {
	id: ID!
	title: String!
}

type Query {
	singlePost(id: ID!): Post
}

schema {
	query: Query
	mutation: Mutation
}
EOF
}

resource "aws_appsync_datasource" "test" {
  api_id           = aws_appsync_graphql_api.test.id
  name             = %[1]q
  service_role_arn = aws_iam_role.test.arn
  type             = "AMAZON_DYNAMODB"

  dynamodb_config {
    table_name = aws_dynamodb_table.test.name
    versioned  = true

    delta_sync_config {
      base_table_ttl        = 60
      delta_sync_table_name = aws_dynamodb_table.test.name
      delta_sync_table_ttl  = 60
    }
  }
}


resource "aws_appsync_resolver" "test" {
  api_id      = aws_appsync_graphql_api.test.id
  field       = "singlePost"
  type        = "Query"
  data_source = aws_appsync_datasource.test.name

  sync_config {
    conflict_detection = "VERSION"
    conflict_handler   = "OPTIMISTIC_CONCURRENCY"
  }

  request_template = <<EOF
{
    "version": "2018-05-29",
    "method": "GET",
    "resourcePath": "/",
    "params":{
        "headers": $utils.http.copyheaders($ctx.request.headers)
    }
}
EOF

  response_template = <<EOF
#if($ctx.result.statusCode == 200)
    $ctx.result.body
#else
    $utils.appendError($ctx.result.body, $ctx.result.statusCode)
#end
EOF
}
`, rName)
}

func testAccResolverConfig_code(rName, code string) string {
	return testAccResolverConfig_base(rName) + fmt.Sprintf(`
resource "aws_appsync_function" "test" {
  api_id      = aws_appsync_graphql_api.test.id
  data_source = aws_appsync_datasource.test.name
  name        = %[1]q
  code        = file("%[2]s")

  runtime {
    name            = "APPSYNC_JS"
    runtime_version = "1.0.0"
  }
}

resource "aws_appsync_resolver" "test" {
  api_id = aws_appsync_graphql_api.test.id
  field  = "singlePost"
  type   = "Query"
  code   = file("%[2]s")
  kind   = "PIPELINE"

  pipeline_config {
    functions = [aws_appsync_function.test.function_id]
  }

  runtime {
    name            = "APPSYNC_JS"
    runtime_version = "1.0.0"
  }
}
`, rName, code)
}
