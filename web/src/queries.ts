import { gql } from "@apollo/client";

export type Webhook = {
  id: string;
  method: string;
  url: string;
  createdAt: string;
};

export const GET_WEBHOOKS = gql`
  query GetWebhooks {
    webhooks {
      nodes {
        id
        url
        method
        createdAt
      }
    }
  }
`;
