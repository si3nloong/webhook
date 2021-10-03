import { gql } from "@apollo/client";

export type Webhook = {
  id: string;
  method: string;
  url: string;
  retries: number;
  createdAt: string;
};

export const GET_WEBHOOKS = gql`
  query GetWebhooks {
    webhooks {
      nodes {
        id
        url
        method
        retries
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
