package namecheap

type Domain struct {
  ID string         `xml:"ID,attr"`
  Name string       `xml:"Name,attr"`
  User string       `xml:"User,attr"`
  Created string    `xml:"Created,attr"`
  Expires string    `xml:"Expires,attr"`
  IsExpired bool    `xml:"IsExpired,attr"`
  IsLocked bool     `xml:"IsLocked,attr"`
  AutoRenew bool    `xml:"AutoRenew,attr"`
  WhoisGuard string `xml:"WhoisGuard,attr"`
}

type ApiResponse struct {
  Status string   `xml:"Status,attr"`
  Command string  `xml:"RequestedCommand"'`
  Domains []Domain `xml:"CommandResponse>DomainGetListResult>Domain"`
}

func (client *NamecheapClient) Domains() ([]Domain, error) {
  resp := ApiResponse{}

  if err := client.get("namecheap.domains.getList", &resp); err != nil {
    return []Domain{}, err
  }
  return resp.Domains, nil
}
