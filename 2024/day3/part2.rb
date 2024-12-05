fname = ARGV[0]
lines = IO.readlines(fname, chomp: true)
count = 0

@flag = true

def mul(n1, n2)
  @flag ? n1 * n2 : 0
end

def flagon()
  @flag = true
end

def dont()
  @flag = false
end

lines.each do |line|
  line.scan(/mul\(\d+,\d+\)|do\(\)|don't\(\)/).each do |expr|
    case expr
    when /don't/ then dont()
    when /do/ then flagon()
    else
      count += eval expr
    end
  end
end

puts count
