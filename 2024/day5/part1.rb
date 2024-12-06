fname = ARGV[0] or raise 'Must supply a filename!'
lines = IO.readlines(fname, chomp: true)
@order = {}

def valid?(nums)
  if nums.empty?
    return true
  end

  n = nums.first
  deps = @order[n]
  rest = nums[1..-1]

  if deps.nil?
    return valid?(rest)
  elsif (deps & rest).empty?
    return valid?(rest)
  else
    return false
  end
end

valid_lines = []

lines.each do |line|
  if line.include?(?|)
    first, second = line.split(?|)
    (@order[second] ||= []) << first
  elsif line.include?(?,)
    data = line.split(?,)
    valid_lines << data if valid?(data)
  end
end

middles = valid_lines.map { |l| l[l.size / 2].to_i }
p middles.reduce(&:+)
