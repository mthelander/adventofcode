fname = ARGV[0]
lines = IO.readlines(fname, chomp: true)
count = 0

def mul(n1, n2)
  n1 * n2
end

lines.each do |line|
  line.scan(/mul\(\d+,\d+\)/).each do |expr|
    count += eval expr
  end
end

puts count
