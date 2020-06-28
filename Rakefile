#!/usr/bin/env rake

require 'yaml'
require 'dotenv'
require 'standalone_migrations'

Dotenv.load

# rake task for database migrations
StandaloneMigrations::Tasks.load_tasks

## uncomment if using mysql with lhm
# require 'lhm'
# # ignore LHM (lhma and lhmn)
# ActiveRecord::SchemaDumper.ignore_tables << /^lhma_/
# ActiveRecord::SchemaDumper.ignore_tables << /^lhmn_/

# namespace :lhm do
#   desc "LHM Clean Up"
#   task cleanup: :environment do
#     Lhm.cleanup(:run)
#   end
# end
