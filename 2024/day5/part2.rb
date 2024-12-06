fname = ARGV[0] or raise 'Must supply a filename!'
lines = IO.readlines(fname, chomp: true)
@order = {}

def find_last_index(nums, vals)
  vals.map { |v| nums.find_index(v) }.compact.max
end

def fix_order!(nums, i: 0)
  return if i >= nums.size

  deps = @order[nums[i]] || []
  violations = deps & nums[i..nums.size]

  unless violations.empty?
    newidx = find_last_index(nums, violations)
    nums.insert(newidx, nums.delete_at(i))
    return fix_order!(nums, i: i)
  end

  fix_order!(nums, i: i+1)
end

result = 0

lines.each do |line|
  if line =~ /[|]/ .. line.empty?
    first, second = line.split(?|)
    (@order[second] ||= []) << first
  else
    nums = line.split(?,)
    fix_nums = nums.clone
    fix_order!(fix_nums)
    result += fix_nums[fix_nums.size / 2].to_i unless nums == fix_nums
  end
end

puts result
