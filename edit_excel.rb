require 'geocoder'
require 'rubygems'
require 'spreadsheet/excel'
require 'roo'
require 'roo-xls'
require 'resolv'

xlsx = Roo::Spreadsheet.open('new.xlsx', extension: :xlsx)

book = Spreadsheet::Workbook.new
sheet1 = book.create_worksheet
Geocoder.configure(timeout: 3)
index = 1

loop do
  begin
    content = xlsx.sheet(0).cell('A', index)
    sheet1.row(index).push content
    index += 1
  end
end

book.write 'edited.xlsx'
