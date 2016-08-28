build=go build

header_set:
	$(build) -o header_set 0_secret.go 0_const_gmail.go 0_service_gmail.go 1_header_set.go

print_header_serial:
	$(build) -o print_header_serial 0_secret.go 0_const_gmail.go 0_service_gmail.go 2_print_header_util.go 2_print_header_serial.go

print_header_parallel_old:
	$(build) -o print_header_parallel_old 0_secret.go 0_const_gmail.go 0_service_gmail.go 2_print_header_util.go 2_print_header_parallel_old.go

print_header_parallel:
	$(build) -o print_header_parallel 0_secret.go 0_const_gmail.go 0_service_gmail.go 2_print_header_util.go 2_print_header_parallel.go

header_meta:
	$(build) -o header_meta 0_secret.go 0_const_gmail.go 0_service_gmail.go 3_header_meta.go

json_read:
	$(build) -o json_read 4_json_read.go

messages:
	$(build) -o 0_secret.go 0_const_gmail.go 0_service_gmail.go 5_messages.go

bq_list:
	$(build) -o bq_list 0_secret.go 0_const_bq.go 0_service_bq.go 6_bq_list.go

json_schema:
	$(build) -o json_schema 7_json_schema.go

bq_upload:
	$(build) -o bq_upload 0_secret.go 0_const_bq.go 0_service_bq.go 8_bq_upload.go

print_header_all:
	$(build) -o print_header_all 0_secret.go 0_service_gmail.go 0_const_gmail.go 9_print_header_all.go

print_message:
	$(build) -o print_message 0_secret.go 0_service_gmail.go 0_const_gmail.go 10_print_message.go

list_messages:
	$(build) -o list_messages 11_list_messages.go

get_messages:
	$(build) -o get_messages 12_get_messages.go
