Lists n8n credentials visible to the API key.

Results are paged internally using cursor pagination until the list is exhausted. Optional `name` and `type` filters are applied after listing (exact match). List is limited to instance owners and admins (`credential:list`); other keys receive HTTP 403.

List items include project `shared` rows. Secret `data` is never included. Live list items do not include `is_managed` or related flags; use `n8n_credential` for those.
