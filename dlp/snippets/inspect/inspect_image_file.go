// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package inspect contains example snippets using the DLP Inspect API.
package inspect

// [START dlp_inspect_file_image]
import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"cloud.google.com/go/dlp/apiv2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

// inspectImageFile inspects an image file at a given filePath, and prints results.
func inspectImageFile(w io.Writer, projectID, filePath string) error {
	// projectID := "my-project-id"
	// filePath := "path/to/image.png"
	// fileType := "IMAGE"
	ctx := context.Background()

	// Initialize client.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	// Gather the resources for the request.
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Create and send the request.
	req := &dlppb.InspectContentRequest{
		Parent: "projects/" + projectID,
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_ByteItem{
				ByteItem: &dlppb.ByteContentItem{
					Type: dlppb.ByteContentItem_IMAGE,
					Data: data,
				},
			},
		},
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes: []*dlppb.InfoType{
				{Name: "PHONE_NUMBER"},
				{Name: "EMAIL_ADDRESS"},
				{Name: "CREDIT_CARD_NUMBER"},
			},
			IncludeQuote: true,
		},
	}
	resp, err := client.InspectContent(ctx, req)
	if err != nil {
		return fmt.Errorf("InspectContent: %v", err)
	}

	// Process the results.
	fmt.Fprintf(w, "Findings: %d\n", len(resp.Result.Findings))
	for _, f := range resp.Result.Findings {
		fmt.Fprintf(w, "\tQoute: %s\n", f.Quote)
		fmt.Fprintf(w, "\tInfo type: %s\n", f.InfoType.Name)
		fmt.Fprintf(w, "\tLikelihood: %s\n", f.Likelihood)
	}
	return nil
}

// [END dlp_inspect_file_image]
