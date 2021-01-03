'use strict';

const { Client } = require('@elastic/elasticsearch');

const client = new Client({ node: 'http://localhost:9200' });

const search = async (query, size = 10) => {
  const result = await client.search({
    index: 'movies',
    size: size,
    body: {
      query: {
        multi_match: {
          query: query,
          fields: [
            'fields.title^4',
            'fields.plot^2',
            'fields.actors',
            'fields.directors',
          ],
        },
      },
    },
  });
  return result.body.hits.hits;
};

const logRecursive = (obj) => {
  console.dir(obj, { depth: null });
};

const main = async () => {
  logRecursive(await client.info());

  if (process.argv.length < 3) {
    throw new Error('query is required');
  }
  const query = process.argv[2];
  const size = process.argv.length === 4 ? process.argv[3] : undefined;
  return search(query, size);
};

main().then(logRecursive).catch(console.error);
