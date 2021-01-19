package ne

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/equinix/ne-go/internal/api"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var testACLTemplate = ACLTemplate{
	Name:        "test",
	Description: "Test ACL",
	MetroCode:   "SV",
	InboundRules: []ACLTemplateInboundRule{
		{
			SrcType:  "SUBNET",
			SeqNo:    1,
			Subnets:  []string{"10.0.0.0/24"},
			Protocol: "TCP",
			SrcPort:  "any",
			DstPort:  "22",
		},
		{
			SrcType:  "DOMAIN",
			SeqNo:    2,
			Subnets:  []string{"216.221.225.13/32"},
			Protocol: "TCP",
			SrcPort:  "any",
			DstPort:  "1024-10000",
		},
	},
}

func TestCreateACLTemplate(t *testing.T) {
	//given
	resp := api.BGPConfigurationCreateResponse{}
	if err := readJSONData("./test-fixtures/ne_acltemplate_post_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	template := testACLTemplate
	reqBody := api.ACLTemplate{}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/ne/v1/device/acl-template", baseURL),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			resp, _ := httpmock.NewJsonResponse(202, resp)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	uuid, err := c.CreateACLTemplate(template)

	//then
	assert.Nil(t, err, "Error is not returned")
	assert.Equal(t, uuid, resp.UUID, "UUID matches")
	verifyACLTemplate(t, template, reqBody)
}

func GetACLTemplates(t *testing.T) {
	//Given
	var respBody api.ACLTemplatesResponse
	if err := readJSONData("./test-fixtures/ne_acltemplates_get_resp.json", &respBody); err != nil {
		assert.Failf(t, "cannot read test response due to %s", err.Error())
	}
	pageSize := respBody.PageSize
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/ne/v1/device/acl-template?size=%d", baseURL, pageSize),
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//When
	c := NewClient(context.Background(), baseURL, testHc)
	c.PageSize = pageSize
	templates, err := c.GetACLTemplates()

	//Then
	assert.Nil(t, err, "Client should not return an error")
	assert.NotNil(t, templates, "Client should return a response")
	assert.Equal(t, len(respBody.Content), len(templates), "Number of objects matches")
	for i := range respBody.Content {
		verifyACLTemplate(t, templates[i], respBody.Content[i])
	}
}

func TestGetACLTemplate(t *testing.T) {
	//given
	resp := api.ACLTemplate{}
	if err := readJSONData("./test-fixtures/ne_acltemplate_get_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	templateID := "db66bf49-b2d8-4e64-8719-d46406b54039"
	testHc := setupMockedClient("GET", fmt.Sprintf("%s/ne/v1/device/acl-template/%s", baseURL, templateID), 200, resp)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	template, err := c.GetACLTemplate(templateID)

	//then
	assert.NotNil(t, template, "Returned template is not nil")
	assert.Nil(t, err, "Error is not returned")
	verifyACLTemplate(t, *template, resp)
}

func TestReplaceACLTemplate(t *testing.T) {
	//given
	templateID := "db66bf49-b2d8-4e64-8719-d46406b54039"
	template := testACLTemplate
	reqBody := api.ACLTemplate{}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/ne/v1/device/acl-template/%s", baseURL, templateID),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			return httpmock.NewStringResponse(204, ""), nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	err := c.ReplaceACLTemplate(templateID, template)

	//then
	assert.Nil(t, err, "Error is not returned")
	verifyACLTemplate(t, template, reqBody)
}

func TestDeleteACLTemplate(t *testing.T) {
	//given
	templateID := "db66bf49-b2d8-4e64-8719-d46406b54039"
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("DELETE", fmt.Sprintf("%s/ne/v1/device/acl-template/%s", baseURL, templateID),
		httpmock.NewStringResponder(204, ""))
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	err := c.DeleteACLTemplate(templateID)

	//then
	assert.Nil(t, err, "Error is not returned")
}

func verifyACLTemplate(t *testing.T, template ACLTemplate, apiTemplate api.ACLTemplate) {
	assert.Equal(t, template.UUID, apiTemplate.UUID, "UUID matches")
	assert.Equal(t, template.Name, apiTemplate.Name, "Name matches")
	assert.Equal(t, template.Description, apiTemplate.Description, "Description matches")
	assert.Equal(t, template.MetroCode, apiTemplate.MetroCode, "MetroCode matches")
	assert.Equal(t, template.DeviceUUID, apiTemplate.VirtualDeviceUUID, "DeviceUUID matches")
	assert.Equal(t, template.DeviceACLStatus, apiTemplate.DeviceACLStatus, "DeviceACLStatus matches")
	assert.Equal(t, len(template.InboundRules), len(apiTemplate.InboundRules), "Number of InboundRules matches")
	for i := range template.InboundRules {
		verifyACLTemplateInboundRule(t, template.InboundRules[i], apiTemplate.InboundRules[i])
	}
}

func verifyACLTemplateInboundRule(t *testing.T, rule ACLTemplateInboundRule, apiRule api.ACLTemplateInboundRule) {
	assert.Equal(t, rule.SeqNo, rule.SeqNo, "SeqNo matches")
	assert.Equal(t, rule.SrcType, rule.SrcType, "SrcType matches")
	assert.Equal(t, rule.FQDN, rule.FQDN, "FQDN matches")
	assert.ElementsMatch(t, rule.Subnets, rule.Subnets, "Subnets matches")
	assert.Equal(t, rule.Protocol, rule.Protocol, "Protocol matches")
	assert.Equal(t, rule.SrcPort, rule.SrcPort, "SrcPort matches")
	assert.Equal(t, rule.DstPort, rule.DstPort, "DstPort matches")
}