# db

Contains file required for building db schema. `schema.rb` is useful when your project already running for a long time and running migration file one by one takes a lot of time, which then you can just load the schema directly which is faster.

The reason I'm using ruby based solution is because there's no golang based solution that can generate schema and allows working on multiple branch with different not yet merged migration without having to rollback and migrate when switching branch.
