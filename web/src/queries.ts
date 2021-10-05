import { gql } from "@apollo/client";

export type Webhook = {
  id: string;
  method: string;
  url: string;
  latestStatusCode: number;
  createdAt: string;
};

export const GET_WEBHOOKS = gql`
  query GetWebhooks {
    webhooks {
      nodes {
        id
        url
        method
        latestStatusCode
        createdAt
      }
    }
  }
`;

export const FIND_WEBHOOK = gql`
  query FindWebhookByID($id: ID!) {
    webhook(id: $id) {
      id
      url
      method
      retries
      createdAt
    }
  }
`;
