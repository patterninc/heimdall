insert into job_statuses
(
    job_status_id,
    job_status_name
)
values
    (1, 'NEW'),
    (2, 'ACCEPTED'),
    (3, 'RUNNING'),
    (4, 'FAILED'),
    (5, 'KILLED'),
    (6, 'SUCCEEDED'),
    (7, 'CANCELLING'),
    (8, 'CANCELLED')
on conflict (job_status_id) do update
set
    job_status_name = excluded.job_status_name;
